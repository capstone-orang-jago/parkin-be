package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/capstone-orang-jago/parkin-be/config"
	"github.com/spf13/cobra"
)

const migrationsDir = "database/migrations"

var migrateCreateCmd = &cobra.Command{
	Use:   "migrate:create [name]",
	Short: "Create a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.ToLower(args[0])
		timestamp := time.Now().Format("20060102150405")
		filename := fmt.Sprintf("%s/%s_%s.sql", migrationsDir, timestamp, name)
		content := "-- +migrate Up\n\n-- +migrate Down\n"

		if err := os.MkdirAll(migrationsDir, os.ModePerm); err != nil {
			fmt.Println("Error creating migrations folder:", err)
			return
		}

		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			fmt.Println("Error creating migration file:", err)
			return
		}

		fmt.Println("Migration created:", filename)
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "migrate:up",
	Short: "Apply all pending up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigrateUp(); err != nil {
			fmt.Fprintf(os.Stderr, "Migration failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "migrate:down",
	Short: "Revert all applied down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMigrateDown(); err != nil {
			fmt.Fprintf(os.Stderr, "Rollback failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCreateCmd)
	rootCmd.AddCommand(migrateUpCmd)
	rootCmd.AddCommand(migrateDownCmd)
}

func runMigrateUp() error {
	conn, err := config.NewDB()
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer conn.Close()

	ctx := context.Background()
	if err := ensureMigrationsTable(ctx, conn); err != nil {
		return err
	}

	applied, err := appliedVersions(ctx, conn)
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("error getting migration files: %w", err)
	}
	sort.Strings(files)

	for _, file := range files {
		version := migrationVersion(file)
		if applied[version] {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}
		upSQL := upSection(string(content))
		if strings.TrimSpace(upSQL) == "" {
			continue
		}

		tx, err := conn.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to start transaction for %s: %w", file, err)
		}

		if _, err := tx.ExecContext(ctx, upSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", file, err)
		}

		fmt.Println("Migrated:", file)
	}

	return nil
}

func runMigrateDown() error {
	conn, err := config.NewDB()
	if err != nil {
		return fmt.Errorf("unable to connect to database: %w", err)
	}
	defer conn.Close()

	ctx := context.Background()
	if err := ensureMigrationsTable(ctx, conn); err != nil {
		return err
	}

	applied, err := appliedVersions(ctx, conn)
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("error getting migration files: %w", err)
	}
	sort.Strings(files)

	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		version := migrationVersion(file)
		if !applied[version] {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}
		downSQL := downSection(string(content))
		if strings.TrimSpace(downSQL) == "" {
			continue
		}

		tx, err := conn.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to start transaction for %s: %w", file, err)
		}

		if _, err := tx.ExecContext(ctx, downSQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to rollback migration %s: %w", file, err)
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM schema_migrations WHERE version = $1`, version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record %s: %w", file, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit rollback %s: %w", file, err)
		}

		fmt.Println("Rolled back:", file)
	}

	return nil
}

func ensureMigrationsTable(ctx context.Context, conn *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`
	if _, err := conn.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}
	return nil
}

func appliedVersions(ctx context.Context, conn *sql.DB) (map[string]bool, error) {
	rows, err := conn.QueryContext(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("failed to read applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan applied migration: %w", err)
		}
		applied[version] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate applied migrations: %w", err)
	}

	return applied, nil
}

func migrationVersion(file string) string {
	return strings.TrimSuffix(filepath.Base(file), ".sql")
}

func upSection(content string) string {
	parts := strings.SplitN(content, "-- +migrate Down", 2)
	return strings.Replace(parts[0], "-- +migrate Up", "", 1)
}

func downSection(content string) string {
	parts := strings.SplitN(content, "-- +migrate Down", 2)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

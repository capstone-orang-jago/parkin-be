package repository

import (
	"context"

	"database/sql"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/sirupsen/logrus"
)

type (
	userRepository struct {
		db *sql.DB
	}
)

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) GetAllUsers(ctx context.Context) (users []*domain.User, err error) {
	query := "SELECT id, name, telephone, email, password, role, saldo, is_verified, created_at, updated_at FROM users"

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		logrus.WithError(err).Error("repo.GetUsers: failed to get users")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &domain.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Telephone, &user.Email, &user.Password, &user.Role, &user.Saldo, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			logrus.WithError(err).Error("repo.GetUsers: failed to scan user")
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repo *userRepository) GetUserByID(ctx context.Context, id string) (user *domain.User, err error) {
	query := `
			SELECT id, name, telephone, email, password, role, saldo, is_verified, created_at, updated_at 
			FROM users 
			WHERE id = $1
			`

	user = &domain.User{}

	err = repo.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Telephone, &user.Email, &user.Password, &user.Role, &user.Saldo, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logrus.WithError(err).Error("repo.GetUserByID: failed to get user by ID")
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (user *domain.User, err error) {
	query := `
			SELECT id, name, telephone, email, password, role, saldo, is_verified, created_at, updated_at 
			FROM users 
			WHERE email = $1
			`

	user = &domain.User{}

	err = repo.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Telephone, &user.Email, &user.Password, &user.Role, &user.Saldo, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logrus.WithError(err).Error("repo.GetUserByEmail: failed to get user by email")
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) CreateUser(ctx context.Context, users *domain.User) (err error) {
	query := `
			INSERT INTO users 
			(id, name, telephone, email, password, role, saldo, is_verified, created_at, updated_at) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
			`
	_, err = repo.db.ExecContext(ctx, query, users.ID, users.Name, users.Telephone, users.Email, users.Password, users.Role, users.Saldo, users.IsVerified, users.CreatedAt, users.UpdatedAt)
	if err != nil {
		logrus.WithError(err).Error("repo.CreateUser: failed to create user")
		return err
	}
	return nil
}

func (repo *userRepository) UpdateUser(ctx context.Context, user *domain.User, id string) error {
	query := `
			UPDATE users 
			SET name = $1, telephone = $2, email = $3, password = $4, role = $5, saldo = $6, is_verified = $7, updated_at = $8 
			WHERE id = $9
			`
	_, err := repo.db.ExecContext(ctx, query, user.Name, user.Telephone, user.Email, user.Password, user.Role, user.Saldo, user.IsVerified, user.UpdatedAt, id)
	if err != nil {
		logrus.WithError(err).Error("repo.UpdateUser: failed to update user")
		return err
	}
	return nil
}

func (repo *userRepository) DeleteUser(ctx context.Context, id string) error {
	query := `
			DELETE FROM users 
			WHERE id = $1
			`
	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		logrus.WithError(err).Error("repo.DeleteUser: failed to delete user")
		return err
	}
	return nil
}

package repository

import (
	"context"

	"database/sql"

	"github.com/capstone-orang-jago/parkin-be/domain"
	"github.com/capstone-orang-jago/parkin-be/helpers"
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
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetUsers: no users found")
		return nil, nil
	} else if err != nil {
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
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetUserByID: no user found for ID")
		return nil, nil
	} else if err != nil {
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
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetUserByEmail: no user found for email")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetUserByEmail: failed to get user by email")
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) GetUserByTelephone(ctx context.Context, telephone string) (user *domain.User, err error) {
	query := `
			SELECT id, name, telephone, email, password, role, saldo, is_verified, created_at, updated_at 
			FROM users 
			WHERE telephone = $1
			`

	user = &domain.User{}

	err = repo.db.QueryRowContext(ctx, query, telephone).Scan(&user.ID, &user.Name, &user.Telephone, &user.Email, &user.Password, &user.Role, &user.Saldo, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetUserByTelephone: no user found for telephone")
		return nil, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetUserByTelephone: failed to get user by telephone")
		return nil, err
	}

	return user, nil
}

func (repo *userRepository) GetUserProfile(ctx context.Context, id, role string) (user *domain.User, err error) {
	user = &domain.User{}
	if role == "petugas" {
		query := `
			SELECT u.id, u.name, u.telephone, u.email, u.role, u.saldo, u.is_verified,
			l.id, l.nama_lokasi, l.alamat, l.is_active
			FROM users u
			LEFT JOIN lokasi_parkir l ON u.id = l.petugas_id
			WHERE u.id = $1 AND u.role = 'petugas'`

		var locID, locNama, locAlamat sql.NullString
		var locActive sql.NullBool

		err = repo.db.QueryRowContext(ctx, query, id).Scan(
			&user.ID, &user.Name, &user.Telephone, &user.Email, &user.Role, &user.Saldo, &user.IsVerified,
			&locID, &locNama, &locAlamat, &locActive,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}

		if locID.Valid {
			user.LokasiParkir = &domain.LokasiParkir{
				ID:         locID.String,
				NamaLokasi: locNama.String,
				Alamat:     locAlamat.String,
				IsActive:   locActive.Bool,
			}
		}

	} else {
		query := `SELECT id, name, telephone, email, role, saldo, is_verified FROM users WHERE id = $1 AND role = 'pengguna'`
		err = repo.db.QueryRowContext(ctx, query, id).Scan(
			&user.ID, &user.Name, &user.Telephone, &user.Email, &user.Role, &user.Saldo, &user.IsVerified,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	}
	return user, nil
}

func (repo *userRepository) CreateUser(ctx context.Context, users *domain.User) (err error) {
	query := `
			INSERT INTO users 
			(name, telephone, email, password, role, saldo, is_verified, created_at, updated_at) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
			`
	_, err = repo.db.ExecContext(ctx, query, users.Name, users.Telephone, users.Email, users.Password, users.Role, users.Saldo, users.IsVerified, users.CreatedAt, users.UpdatedAt)
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

func (repo *userRepository) GetUserSaldo(ctx context.Context, id string) (saldo int64, err error) {
	query := `
			SELECT saldo 
			FROM users 
			WHERE id = $1
			`
	err = repo.db.QueryRowContext(ctx, query, id).Scan(&saldo)
	if err == sql.ErrNoRows {
		logrus.WithError(err).Info("repo.GetUserSaldo: no user found for ID")
		return 0, nil
	} else if err != nil {
		logrus.WithError(err).Error("repo.GetUserSaldo: failed to get user saldo by ID")
		return 0, err
	}

	return saldo, nil
}

func (repo *userRepository) DeductUserSaldo(ctx context.Context, id string, amount int64) error {
	query := `
			UPDATE users 
			SET saldo = saldo - $1, updated_at = NOW() 
			WHERE id = $2
			`

	exec := helpers.GetExecutor(ctx, repo.db)
	_, err := exec.ExecContext(ctx, query, amount, id)
	if err != nil {
		logrus.WithError(err).Error("repo.DeductUserSaldo: failed to deduct user saldo")
		return err
	}
	return nil
}

func (repo *userRepository) AddUserSaldo(ctx context.Context, id string, amount int64) error {
	query := `
			UPDATE users 
			SET saldo = saldo + $1, updated_at = NOW() 
			WHERE id = $2
			`
	exec := helpers.GetExecutor(ctx, repo.db)
	_, err := exec.ExecContext(ctx, query, amount, id)
	if err != nil {
		logrus.WithError(err).Error("repo.AddUserSaldo: failed to add user saldo")
		return err
	}
	return nil
}

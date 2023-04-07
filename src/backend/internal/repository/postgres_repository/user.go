package postgres_repository

import (
	"backend/internal/models"
	"backend/internal/pkg/errors/repositoryErrors"
	"backend/internal/repository"
	"database/sql"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/jmoiron/sqlx"
	"time"
)

type UserPostgres struct {
	UserId    uint64            `db:"user_id"`
	Login     string            `db:"user_login"`
	Password  string            `db:"user_password"`
	Fio       string            `db:"user_fio"`
	DateBirth time.Time         `db:"user_date_birth"`
	Gender    models.UserGender `db:"user_gender"`
	IsAdmin   bool              `db:"user_is_admin"`
}

type UserPostgresRepository struct {
	db *sqlx.DB
}

func NewUserPostgresRepository(db *sqlx.DB) repository.UserRepository {
	return &UserPostgresRepository{db: db}
}

func (i *UserPostgresRepository) Create(user *models.User) error {
	query := `insert into store.users (user_id, user_login, user_password, user_fio,
											 user_date_birth, user_gender, user_is_admin) values
											 ($1, $2, $3, $4, $5, $6, $7);`
	_, err := i.db.Exec(query, user.UserId, user.Login, user.Password, user.Fio, user.DateBirth,
		user.Gender, user.IsAdmin)
	if err != nil {
		return err
	}
	return nil
}

func (i *UserPostgresRepository) Get(login string) (*models.User, error) {
	query := `select * from store.users where user_login = $1`
	userPostgres := &UserPostgres{}

	err := i.db.Get(userPostgres, query, login)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}
	user := &models.User{}
	err = copier.Copy(user, userPostgres)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (i *UserPostgresRepository) GetById(id uint64) (*models.User, error) {
	query := `select * from store.users where user_id = $1`
	userPostgres := &UserPostgres{}

	err := i.db.Get(userPostgres, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repositoryErrors.ObjectDoesNotExists
	} else if err != nil {
		return nil, err
	}
	user := &models.User{}
	err = copier.Copy(user, userPostgres)
	if err != nil {
		return nil, err
	}

	return user, nil
}

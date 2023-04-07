package postgres_repository

import (
	"backend/config"
	"backend/internal/repository"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type PostgresRepositoryFields struct {
	db     *sql.DB
	config config.Config
}

func CreatePostgresRepositoryFields() *PostgresRepositoryFields {
	fields := new(PostgresRepositoryFields)
	err := fields.config.ParseConfig("config.json", "../../../config")
	if err != nil {
		return nil
	}
	fields.db, err = fields.config.Postgres.InitDB()
	if err != nil {
		return nil
	}
	return fields
}

func CreateInstrumentPostgresRepository(fields *PostgresRepositoryFields) repository.InstrumentRepository {
	dbx := sqlx.NewDb(fields.db, "pgx")

	return NewInstrumentPostgresRepository(dbx)
}

func CreateComparisonListPostgresRepository(fields *PostgresRepositoryFields) repository.ComparisonListRepository {
	dbx := sqlx.NewDb(fields.db, "pgx")

	return NewComparisonListPostgresRepository(dbx)
}

func CreateDiscountPostgresRepository(fields *PostgresRepositoryFields) repository.DiscountRepository {
	dbx := sqlx.NewDb(fields.db, "pgx")

	return NewDiscountPostgresRepository(dbx)
}

func CreateUserPostgresRepository(fields *PostgresRepositoryFields) repository.UserRepository {
	dbx := sqlx.NewDb(fields.db, "pgx")

	return NewUserPostgresRepository(dbx)
}

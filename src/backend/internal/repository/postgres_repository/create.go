package postgres_repository

import (
	"backend/config"
	"backend/internal/repository"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type PostgresRepositoryFields struct {
	Db     *sql.DB
	Config config.Config
}

func CreatePostgresRepositoryFields(fileName, filePath string) *PostgresRepositoryFields {
	fields := new(PostgresRepositoryFields)
	err := fields.Config.ParseConfig(fileName, filePath)
	if err != nil {
		return nil
	}
	fields.Db, err = fields.Config.Postgres.InitDB()
	if err != nil {
		return nil
	}
	return fields
}

func CreateInstrumentPostgresRepository(fields *PostgresRepositoryFields) repository.InstrumentRepository {
	dbx := sqlx.NewDb(fields.Db, "pgx")

	return NewInstrumentPostgresRepository(dbx)
}

func CreateComparisonListPostgresRepository(fields *PostgresRepositoryFields) repository.ComparisonListRepository {
	dbx := sqlx.NewDb(fields.Db, "pgx")

	return NewComparisonListPostgresRepository(dbx)
}

func CreateDiscountPostgresRepository(fields *PostgresRepositoryFields) repository.DiscountRepository {
	dbx := sqlx.NewDb(fields.Db, "pgx")

	return NewDiscountPostgresRepository(dbx)
}

func CreateUserPostgresRepository(fields *PostgresRepositoryFields) repository.UserRepository {
	dbx := sqlx.NewDb(fields.Db, "pgx")

	return NewUserPostgresRepository(dbx)
}

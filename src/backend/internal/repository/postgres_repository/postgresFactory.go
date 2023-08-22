package postgres_repository

import (
	"backend/config"
	"backend/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
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

func CreateOrderPostgresRepository(fields *PostgresRepositoryFields) repository.OrderRepository {
	dbx := sqlx.NewDb(fields.Db, "pgx")

	return NewOrderPostgresRepository(dbx)
}

func SetupTestDatabase(migrationFilePath string) (testcontainers.Container, *sql.DB) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	dbContainer, _ := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})

	host, _ := dbContainer.Host(context.Background())
	port, _ := dbContainer.MappedPort(context.Background(), "5432")

	dsnPGConn := fmt.Sprintf("host=%s port=%d user=postgres password=postgres dbname=testdb sslmode=disable", host, port.Int())
	db, err := sql.Open("pgx", dsnPGConn)
	if err != nil {
		return dbContainer, nil
	}

	err = db.Ping()
	if err != nil {
		return dbContainer, nil
	}
	db.SetMaxOpenConns(10)

	text, err := os.ReadFile(migrationFilePath)
	if err != nil {
		return dbContainer, nil
	}

	if _, err := db.Exec(string(text)); err != nil {
		return dbContainer, nil
	}

	return dbContainer, db
}

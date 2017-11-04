package thisapp

import (
	"fmt"
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ProductDB interface {
	GetAll(ctx context.Context, orderBy, orderFrom string) error
}

func NewProductDB(config *Configs) (*PostgresProductDB, error) {
	return newPostgresDB(config)
}


type PostgresProductDB struct {
	connection *sqlx.DB
}

func newPostgresDB(config *Configs) (*PostgresProductDB, error) {
	db, err := sqlx.Connect("postgres", config.Databases.Conn)
	if err != nil {
		return nil, fmt.Errorf("DB open connection error. Error: %s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("DB ping connection error. Error: %s", err.Error())
	}

	database := &PostgresProductDB{
		connection: db,
	}

	return database, nil
}

func (db *PostgresProductDB) GetAll(ctx context.Context, orderBy, orderFrom string) error {
	return nil
}

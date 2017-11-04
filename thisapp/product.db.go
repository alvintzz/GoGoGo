package thisapp

import (
	"fmt"
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ProductDB interface {
	GetAll(ctx context.Context, orderBy, orderFrom string) ([]Product, error)
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

type Product struct {
	ProductID   int    `db:"product_id"`
	ProductName string `db:"product_name"`
}

var getAllQuery = "SELECT product_id, product_name FROM tbl_products ORDER BY %s %s"
func (db *PostgresProductDB) GetAll(ctx context.Context, orderBy, orderFrom string) ([]Product, error) {
	products := []Product{}

	query := fmt.Sprintf(getAllQuery, orderBy, orderFrom)
	err := db.connection.Select(&products, query)
	if err != nil {
		err = fmt.Errorf("Failed to get Product because %s", err.Error())
	}

	return products, err
}

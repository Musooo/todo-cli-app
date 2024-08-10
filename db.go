package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	//"golang.org/x/crypto/bcrypt"
)

type AccountActions interface {
	CreateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
}

type PostgresDb struct {
	db *sql.DB
}

func NewPostgresDb() (*PostgresDb, error) {
	connStr := "user=postgres dbname=test_db password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDb{
		db: db,
	}, nil
}

func (s *PostgresDb) Init() error {
	return s.createAccountTable()
}

func (s *PostgresDb) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		user_name varchar(100),
		encrypted_password varchar(100)
	)`

	_, err := s.db.Exec(query)
	return err
}


package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AccountActions interface {
	CreateAccount(*Account) error
	GetUserByUserName(*string, *string) (*Account, error)
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
	err := s.createAccountTable()
	if err != nil {
		return err
	}

	err = s.createToDoTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDb) createAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		user_name varchar(100),
		password varchar(100)
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresDb) createToDoTable() error {
	query := `create table if not exists todo (
		id INT PRIMARY KEY REFERENCES account(id),
		text VARCHAR(255),
		serial_number SERIAL,
		status bool
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresDb) CreateAccount(acc *Account) error{
	query := `insert into account
	(user_name,password)
	values ($1, $2)`

	tempPass, err := hashPassword(acc.Password)

	if err != nil{
		log.Fatal(err)
	}

	acc.Password = tempPass

	_, err = s.db.Query(
		query,
		acc.UserName,
		acc.Password,
	)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *PostgresDb) IsAccountTaken(userName *string) (bool, error){
	rows, err := s.db.Query("select * from account where user_name = $1", *userName)
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next(){
		return true, nil
	}

	return false, nil
}

func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
    return string(bytes), err
}

func (s *PostgresDb) GetUserByUserName(userName, password *string) (*Logged, error){
	rows, err := s.db.Query("select * from account where user_name = $1", *userName)
	if err != nil{
		log.Fatal(err)
	}
	defer rows.Close()

	var id int
	var passwordGot string
	if err = rows.Scan(&id, &userName, &passwordGot); err != nil {
        log.Fatal(err)
    }
	status := checkPasswordHash(*password, passwordGot)
	if status{
		return NewLogged(*NewAccount(*userName, *password, id)), nil
	}
	log.Fatal("password wrong")
	return nil, nil
}

func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
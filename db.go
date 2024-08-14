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
	GetUserIdByUserName(*string) int
	CreateTodo(*string)
	GetTodos(string) (ToDoArr, error)
	RemoveTodo(int) error
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
	if err := s.db.Ping(); err != nil {
        log.Fatal(err)
    }

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
	query := `
	CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		user_name VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresDb) createToDoTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS todo (
		id SERIAL PRIMARY KEY,
		account_id INT NOT NULL REFERENCES account(id),
		text VARCHAR(255) NOT NULL,
		status BOOLEAN NOT NULL,
		serial_number SERIAL
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

func (s *PostgresDb) CreateTodo(todo ToDo) (int, error){
	query := `insert into todo
	(account_id, text, status)
	values ($1, $2, $3)
	returning id`
	var id int
	err := s.db.QueryRow(
		query,
		todo.UserId,
		todo.Text,
		todo.Status,
	).Scan(&id)

	if err != nil {
		return 0,err
	}
	return id, nil
}

func (s *PostgresDb) IsAccountTaken(userName *string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM account WHERE user_name = $1)", *userName).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	return exists, nil
}


func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 16)
    return string(bytes), err
}

func (s *PostgresDb) GetUserByUserName(userName, password *string) (*Logged, error){
	row := s.db.QueryRow("select * from account where user_name = $1", *userName)

	var id int
	var passwordGot string
	if err := row.Scan(&id, &userName, &passwordGot); err != nil {
        log.Fatal(err)
    }
	status := checkPasswordHash(*password, passwordGot)
	if status{
		return NewLogged(*NewAccount(*userName, *password, id)), nil
	}
	log.Fatal("password wrong")
	return nil, nil
}

func (s *PostgresDb) GetUserIdByUserName(userName *string) int{
	row := s.db.QueryRow("select * from account where user_name = $1", *userName)
	var id int
	var passwordGot string
	if err := row.Scan(&id, &userName, &passwordGot); err != nil {
        log.Fatal(err)
    }
	return id
}

func checkPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func (s *PostgresDb) GetTodos(userId int) (ToDoArr, error) {
	var todos ToDoArr= ToDoArr{}
    query := `SELECT id, account_id, text, status FROM todo WHERE account_id = $1`
    rows, err := s.db.Query(query, userId)
    if err != nil {
        return todos, err
    }
    defer rows.Close()

    for rows.Next() {
        var todo ToDo
        if err := rows.Scan(&todo.ID, &todo.UserId, &todo.Text, &todo.Status); err != nil {
            return todos, err
        }
        todos.ToDos = append(todos.ToDos, todo)
    }

    if err := rows.Err(); err != nil {
        return todos, err
    }

    return todos, nil
}

func (s *PostgresDb) RemoveTodo(todoId int) error{
	query := `DELETE FROM todo WHERE id = $1`
    _, err := s.db.Exec(query, todoId)
    if err != nil {
        return err
    }
	return nil
}
package main

type Account struct {
	ID int `json:"id"`
	UserName string `json:"username"`
	Password string	`json:"password"`
}

type Logged struct {
	UserName string `json:"username"`
	Status bool `json:"status"`
}

type ToDo struct {
	ID int `json:"id"`
	Text string `json:"text"`
}

type ToDoArr struct {
	ToDOArr []ToDo `json:"todos"`
}

type Data struct {
	Accounts []Logged `json:"accounts"`
}

func NewAccount(userName, password string, iD int) *Account{
	return &Account{
		ID: iD,
		UserName: userName,
		Password: password,
	}
}

func NewLogged(user Account) *Logged{
	return &Logged{
		UserName: user.UserName,
		Status: true,
	}
}
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var fileName string = "logs.json"
var fileTodo string = "todos.json"

func main() {
	db, err := NewPostgresDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.db.Close() //i love defer
	
	if err := db.Init(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
	switch os.Args[1] {
		case "login":
			var user *Logged
			user, err = db.GetUserByUserName(&os.Args[2], &os.Args[3])
			if err != nil {
				log.Fatal(err)
			}
			data := getLogs()
			jsonWriting(*user, *data)
			data = getLogs()
			logoutAccs(os.Args[2], *data)
			clearTodosJson()
		case "register":
			var canCreate bool
			canCreate,err = db.IsAccountTaken(&os.Args[2])
			if err != nil{
				log.Fatal(err)
			}
			if !canCreate{
				acc := NewAccount(os.Args[2], os.Args[3], 0)
				db.CreateAccount(acc)
				data := getLogs()
				user := NewLogged(*NewAccount(os.Args[2], os.Args[3], -1))
				jsonWriting(*user, *data)
				data = getLogs()
				logoutAccs(os.Args[2], *data)
				clearTodosJson()
			}else{
				fmt.Println("username already used")
			}
		case "logout":
			data := getLogs()
			logout(os.Args[2], *data) //TODO maybe make it so that you don't need a name, only one acc at a time can be true
			clearTodosJson()
		case "addTodo":
			userName, proced := checkLogged()
			if proced {
				id := db.GetUserIdByUserName(&userName)
				todo := NewTodo(id, os.Args[2]) //maybe more arg
				db.CreateTodo(*todo)
				todoArr := getTodos()
				jsonWTodo(*todo,*todoArr)
			}else {
				fmt.Print("you are not logged in")
			}
		case "list":
			fmt.Print("print all the todos")
		default:
			fmt.Printf("%s is not a recognised command:", os.Args[1])
			os.Exit(1)
	}

}

func getLogs() *Data {
	var data Data

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		data = Data{}
	} else {
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println("Error during file reading:", err)
			return nil
		}

		// dedcoding
		if len(fileContent) == 0 {
            // File is empty, initialize data as an empty object
            data = Data{}
        } else {
            // Decoding
            err = json.Unmarshal(fileContent, &data)
            if err != nil {
                fmt.Println("Error during decoding JSON:", err)
                return nil
            }
        }
	}
	return &data
}

func getTodos() *ToDoArr{
	var toDOArr ToDoArr

	if _, err := os.Stat(fileTodo); os.IsNotExist(err) {
		toDOArr = ToDoArr{}
	} else {
		fileContent, err := os.ReadFile(fileTodo)
		if err != nil {
			fmt.Println("Error during file reading:", err)
			return nil
		}

		// dedcoding
		if len(fileContent) == 0 {
            // File is empty, initialize data as an empty object
            toDOArr = ToDoArr{}
        } else {
            // Decoding
            err = json.Unmarshal(fileContent, &toDOArr)
            if err != nil {
                fmt.Println("Error during decoding JSON:", err)
                return nil
            }
        }
	}
	return &toDOArr
}

func updateJson(data Data) {
	updatedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error during JSON cod:", err)
		return
	}

	err = os.WriteFile(fileName, updatedData, 0644)
	if err != nil {
		fmt.Println("Error during filewriting:", err)
		return
	}

}

func updateJsonTodo(todoArr ToDoArr){
	updatedTodo, err := json.MarshalIndent(todoArr, "", "  ")
	if err != nil {
		fmt.Println("Error during JSON cod:", err)
		return
	}

	err = os.WriteFile(fileTodo, updatedTodo, 0644)
	if err != nil {
		fmt.Println("Error during filewriting:", err)
		return
	}
}

func jsonWTodo(todo ToDo, todoArr ToDoArr){
	todoArr.ToDos = append(todoArr.ToDos, todo)
	updateJsonTodo(todoArr)
}

func jsonWriting(user Logged, data Data) {
	proc, index := isUserInLogs(&user.UserName,&data)
	if proc{
		data.Accounts[index].Status = true 
	}else{
		data.Accounts = append(data.Accounts, user)
	}
	updateJson(data)
}

func isUserInLogs(userName *string, data *Data) (bool, int) {
	for i:=0; i< len(data.Accounts); i++ {
		if data.Accounts[i].UserName == *userName{
			return true,i
		}
	}
	return false, -1
}

func logout(userName string, data Data) {
	for i := 0; i < len(data.Accounts); i++ {
		if data.Accounts[i].UserName == userName {
			data.Accounts[i].Status = false
			updateJson(data)
			break
		}
	}
}


func logoutAccs(userName string, data Data){
	for i:=0; i< len(data.Accounts); i++ {
		if data.Accounts[i].UserName != userName {
			data.Accounts[i].Status = false
			
		}
	}
	updateJson(data)
}

func checkLogged() (string, bool){
	data := getLogs()
	for i:=0; i< len(data.Accounts); i++ {
		if data.Accounts[i].Status{
			return data.Accounts[i].UserName, true
		}
	}
	return "err", false
}

func clearTodosJson(){
	err := os.WriteFile(fileTodo, []byte{}, 0644)
	if err != nil {
		fmt.Println("Error during filewriting:", err)
		return
	}
}
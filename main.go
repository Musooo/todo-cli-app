package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var fileName string = "logs.json"

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
			}else{
				fmt.Println("username already used")
			}
		case "logout":
			data := getLogs()
			logout(os.Args[2], *data)
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

func jsonWriting(user Logged, data Data) {
	data.Accounts = append(data.Accounts, user)
	updateJson(data)
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
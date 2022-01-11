package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	DB_USER 	= "postgres"
	DB_PASSWORD	= "Cuongnguyen2001"
	DB_NAME		= "RestfulApi"
)

type User struct {
	Id 		string `json:"id"`
	Name 	string `json:"name"`
	Gender	string `json:"gender"`
	Email 	string `json:"email"`
	Birth 	string `json:"birth"`
}

type JsonResponse struct {
	Type 	string `json:"type"`
	Data 	[]User `json:"data"`
	Message	string `json:"message"`
}

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		panic(err)
	}

	return db
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/getall", GetAll).Methods("GET")
	router.HandleFunc("/create", CreateUser).Methods("POST")
	router.HandleFunc("/delete/{userid}", DeleteUser).Methods("DELETE")
	router.HandleFunc("/update/{userid}", UpdateUser).Methods("PATCH")
	router.HandleFunc("/find/{userid}", FindUser).Methods("GET")

	fmt.Println("Server at 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func printMessage(message string) {
	fmt.Println(message)
}

// func checkErr(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

func GetAll(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	rows, err := db.Query("SELECT * FROM \"User\".users")
	if err != nil {
		var response = JsonResponse{Type: "Error", Message: "Fail"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var listUser []User

	for rows.Next() {
		var id string
		var name, gender, email, birth string 
		
		err = rows.Scan(&id, &name, &gender, &email, &birth)
		if err != nil {
			var response = JsonResponse{Type: "Error", Message: "Fail"}
			json.NewEncoder(w).Encode(response)
			return
		}

		listUser = append(listUser, User{Id: id, Name: name, Gender: gender, Email: email, Birth: birth})

	}

	var response = JsonResponse{Type: "Success", Data: listUser}
	json.NewEncoder(w).Encode(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	name := r.FormValue("name")
	gender := r.FormValue("gender")
	email := r.FormValue("email")
	birth := r.FormValue("birth")

	var response = JsonResponse{}

	if id == "" ||name == "" || gender == "" || email == "" || birth == "" {
		response = JsonResponse{Type: "Error", Message: "You are missing some parameter"}
	} else {
		db := setupDB()
		printMessage("Insert User to DB")

		var lastInsertID int
		query := fmt.Sprintf("INSERT INTO \"User\".users(id, name, gender, email, birth) VALUES ('{%s}', '{%s}', '{%s}', '{%s}', '{%s}') RETURNING ID;",id, name, gender, email, birth)
		err := db.QueryRow(query).Scan(&lastInsertID)

		if err != nil {
			var response = JsonResponse{Type: "Error", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		response = JsonResponse{Type: "Success", Message: "Creating User is successful"}
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/delete/"):]

	var response = JsonResponse{}

	if id == "" {
		response = JsonResponse{Type: "Error", Message: "You are missing some parameters"}

	} else {
		db := setupDB()
		printMessage("Deleting user from DB")
		query := fmt.Sprintf("DELETE FROM \"User\".users where id = '{%s}'", id)
		_, err := db.Exec(query)
		if err != nil {
			var response = JsonResponse{Type: "Error", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		response = JsonResponse{Type: "Success", Message: "User has been deleted"}

	}
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/update/"):]
	name := r.FormValue("name")
	gender := r.FormValue("gender")
	email := r.FormValue("email")
	birth := r.FormValue("birth")

	var response = JsonResponse{}

	if name == "" || gender == "" || email == "" || birth == "" {
		response = JsonResponse{Type: "Error", Message: "You are missing some parameter"}
	} else {
		db := setupDB()
		printMessage("update User to DB")

		var lastInsertID int
		query := fmt.Sprintf("UPDATE \"User\".users SET name = '{%s}', gender = '{%s}', email = '{%s}', birth = '{%s}' WHERE id = '{%s}' RETURNING ID;", name, gender, email, birth, id)
		err := db.QueryRow(query).Scan(&lastInsertID)

		if err != nil {
			var response = JsonResponse{Type: "Error", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		response = JsonResponse{Type: "Success", Message: "Updating User is successful"}
	}
	json.NewEncoder(w).Encode(response)
}

func FindUser(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	name := r.URL.Path[len("/find/"):]
	var response JsonResponse

	query := fmt.Sprintf("SELECT * FROM \"User\".users WHERE name = '{%s}'", name)
	rows, err := db.Query(query)
	if err != nil {
		response = JsonResponse{Type: "Error", Message: err.Error()}
		return
	}

	var listUser = make([]User, 0)
	for rows.Next() {
		var id string
		var gender, email, birth string 
		
		fmt.Println(name)
		err = rows.Scan(&id, &name, &gender, &email, &birth)
		if err != nil {
			var response = JsonResponse{Type: "Error", Message: err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}

		listUser = append(listUser, User{Id: id, Name: name, Gender: gender, Email: email, Birth: birth})
	}

	response = JsonResponse{Type: "Success", Data: listUser}
	json.NewEncoder(w).Encode(response)
}
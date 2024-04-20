package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var db *sqlx.DB

// var schema = `
// CREATE TABLE person(
// 	id SERIAL,
// 	name text,
// 	email text
// );
// `

type person struct {
	Id    int64
	Name  string
	Email string
}

type users struct {
	Name string
}

func main() {
	var err error
	db, err = sqlx.Connect("postgres", "user=postgres dbname=postgres password=12345678 sslmode = disable")
	if err != nil {
		log.Fatalln(err)
	}
	// create table and insert into table
	// db.MustExec(schema)
	// tx := db.MustBegin()
	// tx.MustExec("INSERT INTO person (name, email) VALUES ($1, $2)", "bhuwan", "bhuwan@gmail.com")
	// tx.Commit()

	people := []person{}
	db.Select(&people, "SELECT * FROM person ")
	for _, all := range people {
		fmt.Printf("users id: %d ,name: %s, email:%s\n", all.Id, all.Name, all.Email)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is backend code for crud"))
	})
	r.Route("/peoples", func(r chi.Router) {
		r.Get("/", getUsers)
		r.Get("/all", getUsersAll)
		r.Post("/", createUser)
		r.Put("/{id}", updateUser)
		r.Delete("/{id}", deleteUser)
	})
	port := "4000"
	log.Printf("Listening on port %s", port)
	http.ListenAndServe(":"+port, r)

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var people []users

	query := "SELECT name FROM person"
	err := db.Select(&people, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(people)
}
func getUsersAll(w http.ResponseWriter, r *http.Request) {
	var people []person

	query := "SELECT * FROM person"
	err := db.Select(&people, query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(people)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var people person
	err := json.NewDecoder(r.Body).Decode(&people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	query := "INSERT INTO  person(name,email) VALUES($1,$2) RETURNING id"
	err = db.QueryRow(query, people.Name, people.Email).Scan(&people.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(people)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := "DELETE FROM person WHERE id = $1"
	_, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("task with ID  %s  deleted successfully", id)))
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var people person
	err := json.NewDecoder(r.Body).Decode(&people)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	query := "UPDATE person SET (name,email) = ($1,$2) WHERE id = $3"
	_, err = db.Exec(query, people.Name, people.Email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("task with ID  %s  Updated successfully", id)))
}

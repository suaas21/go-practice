package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//Book studio
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

//Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var books []Book
var r = mux.NewRouter()

func init() {
	books = append(books, Book{ID: "1", Isbn: "4568", Title: "one Book", Author: &Author{Firstname: "Sayf", Lastname: "Azad"}})
	books = append(books, Book{ID: "2", Isbn: "2569", Title: "Two Book", Author: &Author{Firstname: "Nazim", Lastname: "Uddin"}})
	r.HandleFunc("/books", GetBooks).Methods("GET")
	r.HandleFunc("/books/{id}", GetBook).Methods("GET")
	r.HandleFunc("/books", CreateBooks).Methods("POST")
	r.HandleFunc("/books/{id}", UpdateBooks).Methods("PUT")
	r.HandleFunc("/books/{id}", DeleteBooks).Methods("DELETE")
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get books")
	w.Header().Set("Content-Type", "aplication/json")
	json.NewEncoder(w).Encode(books)
}
func GetBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get book")

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params

	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(Book{})
}
func CreateBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create books")

	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}
func UpdateBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update books")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]

			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(books)

}
func DeleteBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete books")

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
		}
	}
	json.NewEncoder(w).Encode(books)
}
func main() {
	fmt.Println("books server")
	log.Fatal(http.ListenAndServe(":8081", r))

	//first()
	//second()
}

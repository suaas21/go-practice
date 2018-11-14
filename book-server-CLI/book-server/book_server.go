package bookserver

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

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

var (
	empty    = "there is no book"
	notFound = "request book not found"
	edited   = "edited succesfully"
	deleted  = "deleted successfully"
)

var authn map[string]string = map[string]string{"sagor": "azad", "tuhin": "tuhin"}

var books []Book
var loggedIn bool

var r = mux.NewRouter()

func init() {
	books = append(books, Book{ID: "1", Isbn: "4568", Title: "one Book", Author: &Author{Firstname: "Sayf", Lastname: "Azad"}})
	books = append(books, Book{ID: "2", Isbn: "2569", Title: "Two Book", Author: &Author{Firstname: "Nazim", Lastname: "Uddin"}})
}
func CheckAuth(r *http.Request) bool {
	if !loggedIn {
		return true
	}
	encodeUserPass := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	//fmt.Println(encodeUserPass)
	if len(encodeUserPass) != 2 {
		return false
	}
	decodeUserPass, err := base64.StdEncoding.DecodeString(encodeUserPass[1])
	if err != nil {
		return false
	}
	userPass := strings.SplitN(string(decodeUserPass), ":", 2)
	if len(userPass) != 2 {
		return false
	}
	//fmt.Println(userPass)
	temp := 0
	for key, value := range authn {
		if key == userPass[0] && value == userPass[1] {
			temp++
		}
	}
	if temp > 0 {
		return true
	} else {
		return false
	}

}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get books")
	if !CheckAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Need valid username and password")
		return
	}

	w.Header().Set("Content-Type", "aplication/json")
	//w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(books)
	//err = fmt.Errorf("%v", "test err")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error occured: %s", err)
		return
	}
	//data := w.Header().Get("Code")
	//fmt.Println(data)
}
func GetBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get book")
	if !CheckAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Need valid username and password")

		return
	}
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Gets params

	// Loop through books and find one with the id from the params
	for _, item := range books {
		if item.ID == params["id"] {
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("error occured: %s", err)
				return
			}
			return
		}
	}
	log.Println(notFound)
	err := json.NewEncoder(w).Encode(Book{})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error occured: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)

}
func CreateBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("create books")
	if !CheckAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Need valid username and password")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("error occured: %s", err)
		return
	}
	book.ID = strconv.Itoa(rand.Intn(10000000))
	books = append(books, book)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error occured: %s", err)
	}
	w.WriteHeader(http.StatusOK)

}
func UpdateBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update books")
	if !CheckAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Need valid username and password")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			err := json.NewDecoder(r.Body).Decode(&book)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("error occured: %s", err)
			}
			book.ID = params["id"]

			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error occured: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)

}
func DeleteBooks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete books")
	if !CheckAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Need valid username and password")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
		}
	}
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error occured: %s", err)
	}
	w.WriteHeader(http.StatusOK)
}

func RunServer(port string, logIn bool) {
	loggedIn = logIn
	fmt.Println("books server start.....")
	r.HandleFunc("/books", GetBooks).Methods("GET")
	r.HandleFunc("/books/{id}", GetBook).Methods("GET")
	r.HandleFunc("/books", CreateBooks).Methods("POST")
	r.HandleFunc("/books/{id}", UpdateBooks).Methods("PUT")
	r.HandleFunc("/books/{id}", DeleteBooks).Methods("DELETE")

	//log.Fatal(http.ListenAndServe(":8081", r))
	var wait time.Duration
	srv := &http.Server{
		Addr: ":8081",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/books/{id}", GetBook).Methods("GET")
	router.HandleFunc("/books", GetBooks).Methods("GET")
	router.HandleFunc("/books", CreateBooks).Methods("POST")
	router.HandleFunc("/books/{id}", UpdateBooks).Methods("PUT")
	router.HandleFunc("/books/{id}", DeleteBooks).Methods("DELETE")
	return router
}

func TestGetbooks(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req := httptest.NewRequest("GET", "/books/1", nil)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	Router().ServeHTTP(rr, req)
	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

}
func TestGetbook(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req := httptest.NewRequest("GET", "/books", nil)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	Router().ServeHTTP(rr, req)
	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

}

//new book intre
func TestCreateBooks(t *testing.T) {
	//var booktest []Book
	booktest := Book{ID: "1", Isbn: "4568", Title: "56 Book", Author: &Author{Firstname: "Sayf", Lastname: "Azad"}}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(booktest)
	req := httptest.NewRequest("POST", "/books", b)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	Router().ServeHTTP(rr, req)
	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

}
func TestDeleteBooks(t *testing.T) {

	req := httptest.NewRequest("DELETE", "/books/1", nil)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	Router().ServeHTTP(rr, req)
	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())
	//fmt.Println("delete")

}

func TestUpdateBooks(t *testing.T) {
	booktest := Book{ID: "1", Isbn: "4568", Title: "71 Book", Author: &Author{Firstname: "Sayf", Lastname: "Azad"}}
	b := new(bytes.Buffer)
	_ = json.NewEncoder(b).Encode(booktest)
	req := httptest.NewRequest("PUT", "/books/2", b)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	Router().ServeHTTP(rr, req)
	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())
	//fmt.Println("delete")

}

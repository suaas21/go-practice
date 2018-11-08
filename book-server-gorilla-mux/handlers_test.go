package main

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	return router
}

func TestGetbooks(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req := httptest.NewRequest("GET", "/books/2", nil)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	Router().ServeHTTP(rr, req)
	fmt.Println(rr.Code)
	fmt.Println(rr.Body.String())

}

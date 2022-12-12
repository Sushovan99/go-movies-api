package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Movie struct {
	ID       string    `json:"id"`
	ISBN     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Fail struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var movies []Movie
var fail Fail

func getAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content/Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content/Type", "application/json")
	params := mux.Vars(r)

	fail = Fail{
		Status:  "Fail",
		Message: "Not found",
	}

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(fail)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	fail = Fail{
		Status:  "Fail",
		Message: "Not found",
	}

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[:index+1]...)
			w.WriteHeader(204)
			break
		}

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(fail)
	}
}

func updateMovie(w http.ResponseWriter, r *http.Request) {

}

func createMovie(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := mux.NewRouter()
	movies = append(movies, Movie{ID: "1", ISBN: "432877", Title: "Avengers 2012", Director: &Director{FirstName: "Joss", LastName: "Whedon"}})
	movies = append(movies, Movie{ID: "2", ISBN: "432553", Title: "Gladiator 2000", Director: &Director{FirstName: "Ridley", LastName: "Scott"}})

	r.HandleFunc("/movies", getAllMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")
	r.HandleFunc("/movie/{id}", updateMovie).Methods("PATCH")
	r.HandleFunc("/movie", createMovie).Methods("POST")

	fmt.Printf("Starting server at 8080...\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

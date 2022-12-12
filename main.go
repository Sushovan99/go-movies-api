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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fail = Fail{
		Status:  "fail",
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
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	fail = Fail{
		Status:  "fail",
		Message: "Not found",
	}

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			json.NewEncoder(w).Encode(item)
			break
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(fail)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			break
		}
	}

	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	movie.ID = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movies)
}

func main() {
	r := mux.NewRouter()
	movies = append(movies, Movie{ID: "1", ISBN: "432877", Title: "Avengers 2012", Director: &Director{FirstName: "Joss", LastName: "Whedon"}})
	movies = append(movies, Movie{ID: "2", ISBN: "432553", Title: "Gladiator 2000", Director: &Director{FirstName: "Ridley", LastName: "Scott"}})

	r.HandleFunc("/movies", getAllMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PATCH")
	r.HandleFunc("/movies", createMovie).Methods("POST")

	fmt.Printf("Starting server at 8080...\n")
	log.Fatal(http.ListenAndServe(":8080", r))
}

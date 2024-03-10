package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       int       `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director,omitempty"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func seedMovies() {
	movies = append(movies, Movie{ID: 1, Isbn: "423847", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: 2, Isbn: "294679", Title: "Second Movie", Director: &Director{Firstname: "Bob", Lastname: "Smith"}})
	movies = append(movies, Movie{ID: 3, Isbn: "847532", Title: "Bee Movie"})
}

func getAllMoviesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getMovieByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	paramsID, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	idx := slices.IndexFunc(movies, func(m Movie) bool {
		return m.ID == paramsID
	})

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(movies[idx])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func createMovieHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	movie := Movie{}
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	movie.ID = rand.Intn(100000)
	movies = append(movies, movie)
	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func updateMovieByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	paramsID, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	idx := slices.IndexFunc(movies, func(m Movie) bool {
		return m.ID == paramsID
	})

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Delete the old found movie
	movies = append(movies[:idx], movies[idx+1:]...)

	// Create new movie
	newMovie := Movie{}
	err = json.NewDecoder(r.Body).Decode(&newMovie)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newMovie.ID = paramsID
	movies = append(movies, newMovie)
	err = json.NewEncoder(w).Encode(movies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func deleteMovieByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	paramsID, err := strconv.Atoi(params["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	for index, movie := range movies {
		if movie.ID == paramsID {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
}

func main() {
	r := mux.NewRouter()

	seedMovies()

	r.HandleFunc("/movies", getAllMoviesHandler).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovieByIDHandler).Methods("GET")
	r.HandleFunc("/movies", createMovieHandler).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovieByIDHandler).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovieByIDHandler).Methods("DELETE")

	fmt.Println("Server listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}

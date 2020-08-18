package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	sudoku "github.com/pako8128/sudoku_solver"
	"github.com/rs/cors"
)

func solveSudoku(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "application/json")

	var fields [9][9]int
	if err := json.NewDecoder(req.Body).Decode(&fields); err != nil {
		writer.Write([]byte("{\"error\": \"Invalid Input\"}"))
		return
	}

	s := sudoku.Sudoku{
		Fields: fields,
	}

	if err := s.Solve(); err != nil {
		writer.Write([]byte("{\"error\": \"Sudoku Unsolvable\"}"))
		return
	}

	if err := json.NewEncoder(writer).Encode(s.Fields); err != nil {
		log.Fatalf("could not serialize %v", s.Fields)
		writer.Write([]byte("{\"error\": \"Internal Server Error\"}"))
	}
}

func main() {
	log.Println("Setting up Router")
	router := mux.NewRouter()

	router.HandleFunc("/api/solve", solveSudoku).Methods("POST")
	log.Println("[+] POST /api/solve")

	log.Println("Wrapping with CORS defaults")
	handler := cors.Default().Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT environment variable must be defined")
	}

	if err := http.ListenAndServe(":" + port, handler); err != nil {
		log.Fatalf("Could not start server, error: %v", err)
	}
}

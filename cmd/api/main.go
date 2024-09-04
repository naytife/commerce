package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/graph"
)

func main() {
	dbase, err := db.Open("postgres://naytife:naytifekey@localhost:5432/naytifedb?search_path=public&sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbase.Close()
	// initialize the repository
	repo := db.NewRepository(dbase)
	// configure the server
	mux := http.NewServeMux()
	mux.Handle("/", graph.NewPlaygroundHandler("/query"))
	mux.Handle("/query", graph.NewHandler(repo))

	// run the server
	port := ":8080"
	fmt.Fprintf(os.Stdout, "🚀 Server ready at http://localhost%s\n", port)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(port, mux))

}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/petrejonn/naytife/config"
	auth "github.com/petrejonn/naytife/internal"
	"github.com/petrejonn/naytife/internal/db"
	"github.com/petrejonn/naytife/internal/graph"
)

func main() {
	env, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbase, err := db.InitDB(env.DATABASE_URL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbase.Close()
	// initialize the repository
	repo := db.NewRepository(dbase)
	// configure the server
	mux := http.NewServeMux()
	mux.Handle("/", graph.NewPlaygroundHandler("/query"))
	mux.Handle("/query", auth.JWTMiddleware()(graph.NewHandler(repo)))
	// log.Println(env.PORT)
	// log.Println(env.AUTH0_DOMAIN)
	// log.Println(env.DATABASE_URL)
	// log.Println(env.AUTH0_AUDIENCE)
	log.Println(os.Getenv("PORT"))
	log.Println(os.Getenv("AUTH0_DOMAIN"))
	log.Println(os.Getenv("DATABASE_URL"))

	// run the server
	port := "0.0.0.0:8080" // + env.PORT
	fmt.Fprintf(os.Stdout, "🚀 Server ready at http://localhost%s\n", port)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(port, mux))

}

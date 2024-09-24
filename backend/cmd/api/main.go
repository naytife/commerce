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
	"github.com/petrejonn/naytife/internal/middleware"
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
	svr := auth.JWTMiddleware()(graph.NewHandler(repo))
	svr = middleware.ShopIDMiddleware(repo)(svr)
	mux.Handle("/query", svr)

	// run the server
	address := ":" + env.PORT
	fmt.Fprintf(os.Stdout, "🚀 Server ready at port %s\n", address)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(address, mux))

}

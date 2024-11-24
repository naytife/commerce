# The "local" environment 
env "local" {
  url = "postgres://naytife:postgres@localhost:5432/naytifedb?search_path=naytife_schema&sslmode=disable"
  migration {
    dir = "file://internal/db/migrations"
  }
}

# The "prod" environment
env "prod" {
  migration {
    dir = "file://internal/db/migrations"
  }
}

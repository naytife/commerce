# The "local" environment 
env "local" {
  url = "postgres://naytife:naytifekey@localhost:5432/naytifedb?search_path=public&sslmode=disable"
  migration {
    dir = "file://internal/db/migrations"
  }
}

# The "prod" environment
env "prod" {
  url = env("DATABASE_URL") 
  migration {
    dir = "file://internal/db/migrations"
  }
}

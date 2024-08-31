# The "local" environment represents our local testings.
env "local" {
  url = "postgres://naytife:naytifekey@:5432/naytifedb?search_path=public&sslmode=disable"
  migration {
    dir = "file://db/migrations"
  }
}
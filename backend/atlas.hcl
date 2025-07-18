# The "local" environment (traditional PostgreSQL for backward compatibility)
env "local" {
  url = "postgres://naytife:postgres@localhost:5432/naytifedb?sslmode=disable"
  migration {
    dir = "file://internal/db/migrations"
  }
}

# The "local_cnpg" environment (CNPG direct connection for local development)
env "local_cnpg" {
  url = "postgres://naytife:postgres@naytife-postgres-rw.naytife.svc.cluster.local:5432/naytifedb?sslmode=require"
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

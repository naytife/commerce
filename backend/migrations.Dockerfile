# Migration-specific Dockerfile
# Uses the lightweight community Alpine variant with shell utilities
FROM arigaio/atlas:0.35.0-community-alpine

# Copy our configuration and migration files
COPY atlas.hcl /atlas.hcl
COPY internal/db/migrations /migrations

# Set the working directory
WORKDIR /

# Use atlas as the default entrypoint
ENTRYPOINT ["atlas"]

# Use the official PostgreSQL image as the base
FROM postgres:alpine

# Set environment variables for PostgreSQL
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=example
ENV POSTGRES_DB=postgres

# Copy any custom initialization scripts (optional)
# Uncomment if you have SQL files to initialize the database
# COPY init.sql /docker-entrypoint-initdb.d/

# Expose the default PostgreSQL port
EXPOSE 5432

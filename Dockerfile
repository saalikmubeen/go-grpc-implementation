# * Build Stage
FROM golang:1.22-alpine3.19 AS builder

# Set the Current Working Directory inside the container
# meaning all the commands will be executed in this directory
WORKDIR /app

# Copy all the files and source code to this WORKDIR directory
COPY . .

# Build the Go app
RUN go build -o main main.go

# Install curl as golang:1.22-alpine3.19 does not have curl installed
RUN apk --no-cache add curl

# Download and install the migrate tool
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz



# * Run Stage
FROM alpine:3.19
WORKDIR /app
# Copy the Pre-built executable binary file from the previous stage
# to the "." (which represents the WORKDIR)
COPY --from=builder /app/main .

# Copy the migrate tool from the previous stage
# COPY --from=builder /app/migrate.linux-amd64 ./migrate

# Till now we have only copied the binary executable file to the container
# but not the other essential and required files and folders to run 
# the app like app.en, db/migrations, etc
COPY app.env .
COPY ./db/migrations ./db/migrations

COPY run_migrations.sh .
# Make the run_migrations.sh script executable
RUN chmod +x run_migrations.sh


# https://github.com/eficode/wait-for
# Copy the wait-for.sh script to the container
# This script is used to wait for the database or any service to be ready before 
# running another command or service that depends on it.
# Usage: ./wait-for.sh <host>:<port> -t <timeout> -- command args
# Example: ./wait-for.sh postgresDB:5432 -t 30 -- echo "DB is up"
COPY wait-for.sh .
# Make the wait-for.sh script executable
RUN chmod +x wait-for.sh


# TO tell docker that the container listens on this specified port
EXPOSE 8080 

# Default command to run when the container starts
# RUN DB migrations before starting the app
CMD ["/app/main"]

# Main entry point of the docker image
ENTRYPOINT ["/app/run_migrations.sh"]


# When the CMD and ENTRYPOINT are combined, the CMD is the default argument to the ENTRYPOINT.
# When CMD is used together with ENTRYPOINT, it will act just as an additional parameters
# that will be passed into the ENTRYPOINT script .
# It's similar to just running the command in the ENTRYPOINT script:
# ENTRYPOINT ["/app/run_migrations.sh", "/app/main"]
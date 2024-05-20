# Golang + Gin + Postgres + Docker + gRPC + NGINX

gRPC and RESTful HTTP implementation of backend server written in Go.

# Technologies used

- **`Gin`** as HTTP web framework
- **`PostgreSQL`** as database
- **`SQLC`** as code generator for SQL
- **`golang-migrate`** for database migration
- **`gRPC`** Remote procedure call framework.
- **`Docker`** for containerizing the application
- **`NGINX`** as a load balancer and reverse proxy
- **`Protoc`** as protocol buffer compiler



## Setup local development


### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [Evans Cli](https://github.com/ktr0731/evans)


    ```bash
    brew install golang-migrate
    ```


- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [Gomock](https://github.com/golang/mock)

    ``` bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```
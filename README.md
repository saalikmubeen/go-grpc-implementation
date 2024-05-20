# Golang + Gin + Postgres + Docker + gRPC + NGINX

gRPC and RESTful HTTP implementation of backend server written in Go.


### Description
A Go-based based implementation of gRPC with Gin, PostgreSQL, Docker, and NGINX. This project
demonstrates how to build a robust backend service with HTTP and gRPC servers, using Go's Gin
framework for HTTP, a separate gRPC server  and a gRPC gateway to handle HTTP requests under the hood.


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


# Installation

1. Clone project

```
git git@github.com:saalikmubeen/go-grpc-implementation.git
```

## Manual

If you aren't a docker person, [Please learn docker 🥲](https://www.docker.com/why-docker)

cd into root project

```
cd go-grpc-implementation
```

`go mod tidy` to install server dependencies

`Setup required environment variables:`

*In the **`app.env`* file, replace the environment variables with your own.*

*Make sure you have postgresQL  installed*


To start the gRPC and HTTP server:

```
make server
```

**`OR`** ( if you don't have make installed)

```
go run main.go
```


*This will start the gRPC server by default on port  50051 and the HTTP Gateway server on port 8080*


## Docker

**`If you use docker, respect++`**

Running project through docker is a breeze. You don't have to do any setup. Just one docker-compose command and magic

`cd go-grpc-implementation`

### 1. To run the project without load balancing

`run:`

```
make dockerup:
```

**`OR`**

```
docker-compose up --build
```

*This will start the gRPC server by default on port  50051 and the HTTP Gateway server on port 8080*

### 2. To run the project with load balancing

`run:`

```
make loadbalancerup
```

**`OR`**

```
docker-compose -f docker-compose-lb.yml up --build
```

*This will run the 4 instances of our Go server in 4 different containers, each running on a different port and an NGINX load balancer to distribute the load between the 4 instances*
Each instance or container will be running two servers, one for gRPC and as an  HTTP Gateway server. The NGINX will do two things:

- Load balance the incoming HTTP requests between the 4 HTTP Gateway servers running in the 4 instances
- Load balance the incoming gRPC requests between the 4 gRPC servers running in the 4 instances


NGINX Load Balancer opens two ports:

- Port 80, which maps to port 3050 for incoming HTTP requests
- Port 9090, which maps to same port 9090  for incoming gRPC requests


To send HTTP requests to the NGINX load balancer (which will distribute the requests between the 4 HTTP Gateway servers),
just open web browser or Postman and send requests to `http://localhost:3050`

To send gRPC requests to the NGINX load balancer (which will distribute the requests between the 4 gRPC servers),
you can use the Evans CLI tool to send gRPC requests to the NGINX load balancer. Run the following command:

```
evans --port 9090 --host localhost -r repl;
```
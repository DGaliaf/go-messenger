# Go Messenger

Go Messenger is a sample messaging application built using the Go programming language. The application allows to create users, create chat with users, send messages to chat, retrieve messages from chat and get all user`s chats

## Features

- Create user
- Create chat with users
- Send messages to chat
- Retrieve messages from chat
- Get all user`s chats

## Tech

The following technologies were used in the development of this project:

- [Golang](https://go.dev)
- [HttpRouter](https://github.com/julienschmidt/httprouter) - HttpRouter by julienschmidt
- [PostgreSQL](https://www.postgresql.org) - [pgx](https://github.com/jackc/pgx) (PostgreSQL Driver and Toolkit) by jackc
- [Swagger](https://swagger.io) - [swag](https://github.com/swaggo/swag) (Automatically generate RESTful API) by swaggo
- [Docker](https://www.docker.com)

## Installation

Go Messenger requires [Golang](https://go.dev/dl/) v1.20.0+ to run.

- Install and start PostgreSQL
- Change configuration file `/configs/config.local.yml`
- Install docker and start the server.

```sh
cd go-messenger
go run app/cmd/app/main.go
```

## Docker
> Go Messenger is very easy to run through Docker container.

##### Prerequisites
Make sure you have Docker and Docker Compose installed on your machine. You can download Docker Desktop from the official website for your operating system.

##### Environment Variables
Before running the application, you'll need to set some environment variables in the docker-compose.yml file. The following environment variables are used in the provided configuration:

`POSTGRES_HOST`: the hostname or IP address of the Postgres database
`POSTGRES_PORT`: the port number on which Postgres is listening
`POSTGRES_USER`: the username for the Postgres user
`POSTGRES_PASSWORD`: the password for the Postgres user
`POSTGRES_DB`: the name of the Postgres database to use
> You can set these environment variables in the docker-compose.yml file or in your shell environment.

##### Build and Run the Application
Change into the root directory of the cloned repository and run the following command to build and start the application:

```ssh
docker-compose up --build
```

> This command will build the Go application using the Dockerfile and start two containers: one for the HTTP server and one for Postgres. The --build flag ensures that the latest version of the application is built before starting the containers.

Once the application is running, you should be able to access it in your web browser by navigating to http://localhost:9000. The HTTP server is running on port 9000, which is mapped to the container's port 9000 in the docker-compose.yml file.

##### Stop the Application
To stop the application and remove the containers, press Ctrl-C in the terminal where docker-compose up is running. Alternatively, you can run the following command in the same directory as the docker-compose.yml file:
```ssh
docker-compose down
```
This command will stop and remove the containers created by the docker-compose up command.

### Swagger Documentation
You can view the Swagger documentation for the API by navigating to `http://localhost:9000/swagger` in your web browser. This will display a UI that shows all the available endpoints, along with details on the request and response formats.
> The Swagger UI provides a convenient way to explore the API and test its endpoints. You can try out each endpoint by clicking on it and entering the required parameters. The UI will display the response in real-time, making it easy to see the results of your requests.

## License

##### MIT

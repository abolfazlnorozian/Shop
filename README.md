# Shop

Online Shop API: A GoLang Web Application with MongoDB, Gin, and JWT

## Project Structure
- [Database](./database/README.md)
- [Entities](./entities/README.md)
- [Auth](./auth/README.md)
- [Helpers](./helpers/README.md)
- [Response](./response/README.md)
- [Router](./router/README.md)
- [Services](./services/README.md)
- [Upload](./upload/README.md)

`go.mod` and `go.sum`: These files are used for managing dependencies with Go modules.

`main.go`: This is the entry point of the application.

`template.env`: This file stores configuration variables or environment settings for the project.

- `database`: Contains files related to the database.
- `entities`: Contains files defining the entity models.
- `auth`: Contains auth files for authentication and authorization.
- `public`: Contains public files, such as images.
- `helpers`: Contains helpers files.
- `response`: Contains files related to response handling.
- `router`: Contains the router file for defining API endpoints.
- `services`: Contains service files for business logic.
- `upload`: Contains files related to file uploading.



## Requirements

go 1.17

MongoDB 4.4.13 

## Set up environment variables

Rename the template.env file to .env

Open the .env file and provide the required configuration variables such as:

`MONGO_URL=`:localhost for MongoDB

`PORT=`:Database connection port number

`DBNAME=`:Insert the desired database name




If all tests pass, start the HTTP service with:

```bash
go run main.go
```
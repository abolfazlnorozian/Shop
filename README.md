# Shop

Online Shop API: A GoLang Web Application with MongoDB, Gin, and JWT

## Project Structure

`database`: This directory contains files related to the database, including the database configuration (`database.go`), environment variables (`env.go`), and functions for working with collections (`getcollection.go`).

`auth`: This directory contains `auth` files responsible for authentication and authorization. It includes files such as `adminToken.go`, `authMiddleware.go`, and `userToken.go`.

`entities`: The` entities` directory houses the entity models for the online shop. It includes files such as `admin.go`, `carts.go`, `category.go`, `counter.go`, `order.go`, `product.go`, `upload.go`, and `user.go`.

`helpers`: The helpers directory contains files that are A series of auxiliary functions to perform special commands, such as code generation or semi-random user, etc. It includes files like `generateId.go` and `generateCode` and `generateUsername.go`and `user.go`.

`response`‍: This directory includes files related to handling API responses, such as `response.go`.

`router`: The router directory contains the router file (`router.go`) responsible for defining API endpoints and routing requests.

`services`: This directory holds the service files that encapsulate the business logic of the online shop. It includes files like `admin.go`, `carts.go`, `category.go`,` order.go`, `product.go`, and `user.go`.

`upload`: The upload directory consists of files related to file uploading functionality, such as `upload.go`.

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
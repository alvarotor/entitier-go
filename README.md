# Entitier-Go Repository

## Overview

This repository contains a collection of Go files that provide a robust structure for building entities and services for a Go application. It implements a generic repository and service pattern, allowing for flexible and reusable data access and business logic layers.

## Directories and Files

### models/

The models directory contains Go files that define the data structures used by the application. Right now error messages.

#### errors.go

The errors.go file defines a set of custom error types that can be used throughout the application. This allows for more specific error handling and reporting.

### repositories/

The repositories directory contains Go files that define the data access layer of the application.

#### generic-repo.go

The generic-repo.go file implements a generic repository using Go's generics feature. It provides the following methods:

- `Create`: Adds a new entity to the database.
- `GetAll`: Retrieves all entities of a specific type.
- `Get`: Retrieves a single entity by ID, with an option to preload related data.
- `Update`: Modifies an existing entity.
- `Delete`: Removes an entity, with an option for soft or hard deletion.

This implementation uses GORM as the ORM (Object-Relational Mapping) library to interact with the database.

#### interface-generic-repo.go

The interface-generic-repo.go file defines the `IGenericRepo` interface, which specifies the methods that any repository implementation should provide. This allows for easy swapping of repository implementations if needed.

### controllers/

The controllers directory contains Go files that define the controllers of the application.

#### generic-controller.go

The generic-controller.go file implements a generic controller that works with the generic repository and service. It provides a layer of abstraction between the repository, service and the application's HTTP handlers. The controller includes methods that correspond to the CRUD operations:

- `GetAll`: Retrieves all entities of a specific type.
- `Get`: Retrieves a single entity by ID.
- `Create`: Creates a new entity.
- `Delete`: Removes an entity.
- `Update`: Modifies an existing entity.

<!-- ### services/

The services directory contains Go files that define the business logic of the application.

#### generic-service.go

The generic-service.go file implements a generic service that works with the generic repository. It provides a layer of abstraction between the repository and the application's business logic. The service includes methods that correspond to the repository operations:

- `GetAll`: Retrieves all entities of a specific type.
- `Get`: Retrieves a single entity by ID.
- `Create`: Creates a new entity.
- `Delete`: Removes an entity.
- `Update`: Modifies an existing entity. -->

### middleware

The middleware directory contains Go files that define the middlewares of the application. Such as authorization, validation, etc.

#### middleware.go

The middleware.go file implements a set of middlewares that can be used to protect routes and validate requests. For example, the `IDValidator` middleware is used to validate ID parameters in routes, ensuring that they are of the correct type and within a valid range.

### utils

The utils directory contains Go files that define the utils of the application. Such as helpers, constants, etc.

#### utils.go

The utils.go file implements a set of utils that can be used throughout the application. For instance, the `GetIDParam` function retrieves the ID parameter from a gin context, handling various cases like missing or invalid ID values. Additionally, the `ConvertToGenericID` function provides a way to convert an `interface{}` type to a generic ID type (string or uint), allowing for flexible ID handling across different entities.

## Key Features

1. **Generic Implementation**: Both the repository and service are implemented using Go's generics, allowing them to work with various entity types.

2. **Flexibility**: The generic approach allows for easy extension to new entity types without duplicating code.

3. **Separation of Concerns**: The repository handles data access, while the service layer manages business logic, promoting a clean architecture.

4. **Error Handling**: Custom error types are used to provide more meaningful error messages and easier error handling.

5. **GORM Integration**: The repository uses GORM, a popular Go ORM, for database operations, providing powerful querying capabilities.

## Usage

To use this structure in your project:

1. Define your entity models in the `models/` directory of your project.
<!-- 2. Create instances of `GenericRepository` and `GenericService` for each of your entity types. -->
2. Create instances of `GenericRepository` for each of your entity types.
3. Use these instances in your application logic to perform CRUD operations on your entities.

Example:

```go
// Assuming you have a User entity
db := // your GORM database instance
userRepo := repositories.NewGenericRepository[User, uint](db)
// userService := services.NewGenericService[User, uint](userRepo)

// Now you can use userService to perform operations on User entities
users, err := userService.GetAll()
```

You can use it with controllers directly too in your project:

```go
package main

import (
    "github.com/alvarotor/entitier-go/controllers"
    "github.com/alvarotor/entitier-go/logger"
    "github.com/alvarotor/entitier-go/middleware"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func main() {
    r := gin.Default()

    // Create a controller instance
    db := // ...initialize DB GORM connection
    log := logger.NewLogger() // Assume you have a logger package
    userController := controllers.NewGenericController[User, uint](log, db)

    // Example route using the IDValidator middleware
    r.GET("/users", userController.GetAll)
    r.GET("/users/:id", middleware.IDValidator[uint](), userController.Get)

    r.Run()
}
```

## Test

Test have been done with `mockery` and installed with `sudo apt install mockery` (linux). Command: `mockery --all --with-expecter`.
To test fully coverage:

```sh
go test -v -coverprofile cover.out ./...
go tool cover -html cover.out -o cover.html
```

And open the `cover.html` file in your browser.

## Contributions

Contributions to improve the generic repository and service implementations or to add new features are welcome. Please submit a pull request or open an issue to discuss proposed changes.

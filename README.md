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

### services/

The services directory contains Go files that define the business logic of the application.

#### generic-service.go

The generic-service.go file implements a generic service that works with the generic repository. It provides a layer of abstraction between the repository and the application's business logic. The service includes methods that correspond to the repository operations:

- `GetAll`: Retrieves all entities of a specific type.
- `Get`: Retrieves a single entity by ID.
- `Create`: Creates a new entity.
- `Delete`: Removes an entity.
- `Update`: Modifies an existing entity.

## Key Features

1. **Generic Implementation**: Both the repository and service are implemented using Go's generics, allowing them to work with various entity types.

2. **Flexibility**: The generic approach allows for easy extension to new entity types without duplicating code.

3. **Separation of Concerns**: The repository handles data access, while the service layer manages business logic, promoting a clean architecture.

4. **Error Handling**: Custom error types are used to provide more meaningful error messages and easier error handling.

5. **GORM Integration**: The repository uses GORM, a popular Go ORM, for database operations, providing powerful querying capabilities.

## Usage

To use this structure in your project:

1. Define your entity models in the `models/` directory.
2. Create instances of `GenericRepository` and `GenericService` for each of your entity types.
3. Use these instances in your application logic to perform CRUD operations on your entities.

Example:

```go
// Assuming you have a User entity
db := // your GORM database instance
userRepo := repositories.NewGenericRepository[User, uint](db)
userService := services.NewGenericService[User, uint](userRepo)

// Now you can use userService to perform operations on User entities
users, err := userService.GetAll()
```

## Contributions

Contributions to improve the generic repository and service implementations or to add new features are welcome. Please submit a pull request or open an issue to discuss proposed changes.

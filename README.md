# Product Management

Product Management is a small project with basic CRUD product, category, order_items, register user, login, create order, authentication, authorization.

## Technologies used:
1. Database: CockroachDB
2. Cache: Redis Cache
3. Logger: Honeycomb, Zap
4. Web framework: Gin
5. Security: JWT Token

## API Document
Run application and click this link
[Documentation](http://localhost:8080/api/v1/swagger/index.html#/)

# Project Architecture
```
📦 
├─ application.go
├─ application
├─ docs
│  ├─ docs.go
│  ├─ swagger.json
│  └─ swagger.yaml
├─ domain
│  ├─ entity
│  └─ repository
├─ go.mod
├─ go.sum
├─ infrastructure
│  ├─ config
│  ├─ controllers
│  │  ├─ handlers
│  │  ├─ middleware
│  │  └─ payload
│  ├─ implementations
│  │  ├─ cache
│  │  ├─ categories
│  │  ├─ loggers
│  │  │  ├─ base.go
│  │  │  ├─ honeycomb
│  │  │  ├─ slack
│  │  │  └─ zap
│  │  ├─ orders
│  │  ├─ products
│  │  ├─ user_roles
│  │  └─ users
│  ├─ jobs
│  ├─ mapper
│  ├─ persistences
│  │  └─ base
│  │     ├─ base.go
│  │     ├─ db
│  │     ├─ logger
│  │     ├─ redis
│  │     └─ supabase
│  └─ routes
└─ utils
```

## Structure

### Root Level

- **application.go**: Main file that starts the application.

### application

- Contains configuration and application startup files.

### docs

- **docs.go**: Describes API documentation parameters and configuration.
- **swagger.json**: JSON file describing the Swagger API.
- **swagger.yaml**: YAML file describing the Swagger API.

### domain

- Defines core components of the application.
    - **entity**: Contains entity definitions in the application domain.
    - **repository**: Contains repository interfaces for entities.

### infrastructure

- Contains infrastructure components and external integrations.
    - **config**: Manages application configuration.
    - **controllers**: Contains controllers that handle HTTP requests.
        - **handlers**: Handles specific requests.
        - **middleware**: Contains application middleware.
        - **payload**: Contains payload definitions for the API.
    - **implementations**: Contains specific implementations for interfaces.
        - **cache**: Implements caching mechanisms.
        - **categories**: Handles business logic related to categories.
        - **loggers**: Implements loggers.
            - **base.go**: Basic logger configuration.
            - **honeycomb**: Implements Honeycomb logger.
            - **slack**: Implements Slack logger.
            - **zap**: Implements Zap logger.
        - **orders**: Handles business logic related to orders.
        - **products**: Handles business logic related to products.
        - **user_roles**: Handles business logic related to user roles.
        - **users**: Handles business logic related to users.
    - **jobs**: Contains background jobs.
    - **mapper**: Contains mappers for data transformations.
    - **persistences**: Contains data storage implementations.
        - **base**: Basic data storage configuration.
            - **base.go**: Basic configuration for data storage.
            - **db**: Database storage implementation.
            - **logger**: Manages loggers within data storage.
            - **redis**: Implements Redis storage.
            - **supabase**: Implements Supabase storage.
    - **routes**: Defines routes for the application.

### utils

- Contains utilities and general support tools for the application.

### go.mod

- File that declares the project's modules and dependencies.

### go.sum

- File that locks the versions of the project's dependencies.

# Installation



```powershell
git clone https://github.com/whisssp/pm_be.git
```

## Run application

1. install package
```powershell
go get .
```

2. run application
```powershell
go run application.main
```
# Product Management

Product Management is a small project with basic CRUD product, category, order_items, register user, login, create order, authentication, authorization.

## Technologies used:
1. Database: CockroachDB
2. Cache: Redis Cache
3. Logger: Honeycomb, Zap
4. Web framework: Gin, Gorm
5. Security: JWT Token

## API Document
Run application and click this link
[Documentation](http://localhost:8080/api/v1/swagger/index.html#/)

# Project Architecture
- This project is using DDD *(Domain-Driven Design)* Architecture 

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

- Contains business logic and use case implementations.

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
    - **implementations/**: Contains implementations of various interfaces.
      - **cache**: Implementations for caching mechanisms.
      - 
      - **categories**: Implementations related to category management.
      - **loggers/**: Logger implementations.
          - **base.go**: Base logger implementation.
          - **honeycomb**: Honeycomb logger implementation.
          - **slack**: Slack logger implementation.
          - **zap**: Zap logger implementation.
      - **orders**: Implementations related to order management.
      - **products**: Implementations related to product management.
      - **user_roles**: Implementations related to user role management.
      - **users**: Implementations related to user management.
    - **jobs**: Contains background jobs.
    - **mapper**: Contains mappers for data transformations.
    - **persistences**: Contains data storage implementations.
        - **base.go**: Initialize basic data storage configuration.
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

# Clone project



```powershell
git clone https://github.com/whisssp/pm_be.git
```

## Run application

1. Go to the project directory
```powershell
cd <path where the project was cloned>/pm_be
```
2. install package
```powershell
go get .
```
3. run application
```powershell
go run application.main
```
# User-Manager

User Manager is a simple application built in Go. It allows users management with basic CRUD operations.

## Features

- **Create**: Add new users with a name and age.
- **Read**: View a specific user by ID.
- **Update**: Modify existing users.
- **Delete**: Remove users.

## Technologies Used

- **Go**: Programming language used for development.
- **PostgreSQL**: Database management system for storing tasks.
- **Gorilla Mux**: HTTP router used for handling routing and middleware.
- **Swagger**: API documentation tool.
- **Testing**: Unit tests written using the standard Go testing package.


## Getting Started

### Prerequisites

To run this project, you need to have the following software installed on your system:

- Go
- PostgreSQL
- swag

### Installation

1. **Clone the repository**:

2. **Set up the database**:

Ensure PostgreSQL is running.
Update configuration in config/config.go.

3. **Build and Run**:

Build and run the application:
```bash
make build
./bin/main
```

4. **Generate Swagger**:

```bash
make swagger
```

5. **Access The API Documentation**:
Open your browser and go to http://localhost:8080/swagger/index.html

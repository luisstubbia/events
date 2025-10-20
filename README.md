# Event Service

A robust RESTful API service for managing events, built with Go and PostgreSQL. Features OpenAPI documentation, comprehensive error handling, and production-ready configuration.

## 🚀 Features

- **RESTful API** - Full CRUD operations for events management
- **OpenAPI Documentation** - Interactive Swagger UI for API exploration
- **PostgreSQL Integration** - Efficient data storage with connection pooling
- **Input Validation** - Comprehensive request validation and error handling
- **CORS Support** - Cross-origin resource sharing enabled
- **Graceful Shutdown** - Proper cleanup on application termination
- **Docker Support** - Containerized deployment ready

## 📋 Prerequisites

- **Go 1.21+** - [Download here](https://golang.org/dl/)
- **PostgreSQL 12+** - [Download here](https://www.postgresql.org/download/)
- **Git** - For version control

### Optional
- **Docker & Docker Compose** - For containerized deployment
- **curl or Postman** - For API testing

## 🛠️ Installation & Setup

### 1. Clone the Repository

```
bash 
git clone <repository-url> 
cd events
```

### 2. Set Up PostgreSQL
**Option A: Using Docker (Recommended)**
```
bash
Start PostgreSQL with Docker
docker run --name event-postgres \
  -e POSTGRES_DB=eventsdb \
  -e POSTGRES_USER=eventdbuser \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  -d postgres:15`
```

Verify container is running
```
docker ps
```

**Option B: Local PostgreSQL Installation**

1. Install PostgreSQL from [official website](https://www.postgresql.org/download/)
2. Start PostgreSQL service:
```
bash
# On macOS with Homebrew
brew services start postgresql

# On Ubuntu
sudo systemctl start postgresql

# On Windows
net start postgresql
```

3. Create database and user:

```
sql
CREATE DATABASE eventsdb;
CREATE USER eventdbuser WITH PASSWORD 'password';
GRANT ALL PRIVILEGES ON DATABASE eventsdb TO eventdbuser;
```

### 3. Database Schema Setup
```
bash
# Run the setup script
chmod +x setup_db.sh
./setup_db.sh

# Or manually execute the SQL
psql -h localhost -U postgres -d eventsdb -f setup_database.sql
```

### 4. Install Go Dependencies
```
bash
# Install Swag CLI for OpenAPI documentation
go install github.com/swaggo/swag/cmd/swag@latest

# Add GOPATH to your PATH (add to your ~/.zshrc or ~/.bashrc)
export PATH=$PATH:$(go env GOPATH)/bin

# Generate Swagger documentation
swag init -g cmd/server/main.go

# Install application dependencies
go mod tidy
```

## 🏃‍♂️ Running the Service
### Development Mode
```
bash
# Run directly with Go
go run cmd/server/main.go

# Or build and run
go build -o events cmd/server/main.go
./events
```
Using Environment Variables
```
bash
# Custom configuration
DB_HOST=localhost \
DB_PORT=5432 \
DB_USER=eventdbuser \
DB_PASSWORD=password \
DB_NAME=eventsdb \
DB_SSLMODE=disable \
go run main.go
```

### Production Mode
```
bash
# Build for production
go build -ldflags="-s -w" -o events cmd/server/main.go

# Run with production settings
./events
```

## 📊 API Endpoints
### Events Management
| Method | Endpoint | Description | Status Codes |
|--------|----------|-------------|--------------|
| **POST** | `/events` | Create a new event | 201, 400, 500 |
| **GET** | `/events` | List all events (ordered by start_time) | 200, 500 |
| **GET** | `/events/{id}` | Get event by UUID | 200, 400, 404, 500 |

### System Endpoints

| Method | Endpoint | Description | Status Codes |
|--------|----------|-------------|--------------|
| **GET** | `/health` | Health check | 200, 503 |
| **GET** | `/swagger/*` | Swagger UI | 200 |
| **GET** | `/` | Redirects to Swagger UI | 302 |

## Request/Response Examples
### Create Event
```
bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "description": "Weekly team sync",
    "start_time": "2024-01-15T10:00:00Z",
    "end_time": "2024-01-15T11:00:00Z"
  }'
  ```
### List Events
```
bash
curl http://localhost:8080/events
```

### Get Event by ID
```
bash
curl http://localhost:8080/events/550e8400-e29b-41d4-a716-446655440000
```

### Health Check
```
bash
curl http://localhost:8080/health
```

## 📚 API Documentation
### Swagger UI
Once the service is running, access the interactive API documentation at:
```
text
http://localhost:8080/swagger/index.html
```
### OpenAPI Specification
The raw OpenAPI specification is available at:
```
text
http://localhost:8080/swagger/doc.json
```

## 🐳 Docker Deployment
### Prerequisites
- Docker Engine 20.10+
- Docker Compose 2.0+

### Quick Start with Docker Compose
1. Clone and navigate to the project:
```
bash
git clone <repository-url>
cd events
```
2. Start the complete stack:
```
bash
docker-compose up -d
```
3. View logs:
```
bash
docker-compose logs -f events
```
4. Stop the stack:
```
bash
docker-compose down
```
### Manual Docker Build
```
bash
# Build the Docker image
docker build -t events .

# Run the container
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=eventdbuser \
  -e DB_PASSWORD=password \
  -e DB_NAME=eventsdb \
  events
  ```

### Docker Compose Configuration
The docker-compose.yml file defines:

- PostgreSQL service with persistent volume
- Event Service application
- Network for inter-service communication
- Health checks for both services

## ⚙️ Configuration
### Environment Variables
| Variable | Default | Description |
|--|----|-------|
|DB_HOST| localhost |PostgreSQL host|
|DB_PORT|5432|PostgreSQL port|
|DB_USER| eventdbuser |Database user|
|DB_PASSWORD|password|Database password|
|DB_NAME|eventsdb|Database name|
|DB_SSLMODE|disable|SSL mode for database connection|

### Database Configuration
The service uses connection pooling with the following settings:

- Max Open Connections: 25
- Max Idle Connections: 25
- Connection Lifetime: 5 minutes
- Query Timeout: 5 seconds

## 🧪 Testing
Manual Testing with curl
```
bash
# Health check
curl http://localhost:8080/health

# Create an event
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Event",
    "description": "Test Description",
    "start_time": "2024-01-20T09:00:00Z",
    "end_time": "2024-01-20T10:00:00Z"
  }'

# List events
curl http://localhost:8080/events

# Get specific event (replace with actual UUID)
curl http://localhost:8080/events/$(curl -s http://localhost:8080/events | jq -r '.[0].id')
```

### Using Swagger UI
1. Start the service
2. Open http://localhost:8080/swagger/index.html
3. Use the "Try it out" feature for each endpoint
4. Execute requests directly from the browser

<img src="/scripts/swagger.png"/>

## 🔧 Development
Project Structure
```
text
event-api/
├── src/
├──── server/
├────── main.go           # Main application entry point
├── go.mod                # Go module dependencies
├── scripts/
├──── schema.sql          # Database schema and sample data
├──── 01-init.sql         # Database setup script
├── docker-compose.yml    # Docker Compose configuration
├── docs/                 # Auto-generated Swagger documentation
├── internal/             # Application source. Hidden for import
├──── domain/             # Application bussines rules implementation.
├──── entrypoint/         # Application entrypoints.
├── pkg/                  # External dependencies like db or web framework.
└── docs
```

## Access Points After Startup:

- 🚀 Application: http://localhost:8080
- 📚 API Docs: http://localhost:8080/swagger/index.html
- ❤️ Health Check: http://localhost:8080/health
- 🗄️ API Endpoint: http://localhost:8080/events

The service will start on port 8080 by default. Check the console output for exact URLs and any startup messages.

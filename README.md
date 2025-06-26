# GoFamTree - Family Tree API

A RESTful API for managing family trees built with Go, GORM, and PostgreSQL.

## Features

- 🔐 Admin authentication with password hashing
- 🏠 House management (family groups)
- 👥 Person management with personal details
- 🔗 Relationship management (parent, spouse, sibling)
- 🌳 Family tree visualization endpoint
- 📊 Full CRUD operations for all entities
- 🛡️ Data validation and constraint enforcement

## Tech Stack

- **Language:** Go 1.21+
- **Database:** PostgreSQL 17
- **ORM:** GORM
- **Architecture:** Clean REST API

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 17
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd gofamtree
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up the database:
```sql
CREATE DATABASE gofamtree_new;
```

4. Configure environment (optional):
```bash
export DATABASE_URL="host=localhost user=postgres dbname=gofamtree_new port=5432 sslmode=disable"
export PORT=8080
```

5. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Admin Authentication

#### Register Admin
```http
POST /admin/register
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

#### Login Admin
```http
POST /admin/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

### House Management

#### Create House
```http
POST /houses
Content-Type: application/json

{
  "name": "Smith Family",
  "created_by": 1
}
```

#### Get All Houses
```http
GET /houses
```

#### Get House by ID
```http
GET /houses/1
```

#### Update House
```http
PUT /houses/1
Content-Type: application/json

{
  "name": "Updated Smith Family"
}
```

#### Delete House
```http
DELETE /houses/1
```

### Person Management

#### Create Person
```http
POST /persons
Content-Type: application/json

{
  "house_id": 1,
  "name": "John Smith",
  "contact": "john@example.com",
  "description": "Father of the family",
  "gender": "male",
  "dob": "1980-01-15"
}
```

#### Get All Persons
```http
GET /persons
# Filter by house:
GET /persons?house_id=1
```

#### Get Person by ID
```http
GET /persons/1
```

#### Update Person
```http
PUT /persons/1
Content-Type: application/json

{
  "name": "John Smith Jr.",
  "contact": "johnsmith@example.com",
  "description": "Updated description",
  "gender": "male",
  "dob": "1980-01-15"
}
```

#### Delete Person
```http
DELETE /persons/1
```

### Relationship Management

#### Create Relationship
```http
POST /relations
Content-Type: application/json

{
  "house_id": 1,
  "person_id": 1,
  "related_to_id": 2,
  "relation_type": "parent"
}
```

Supported relation types: `parent`, `spouse`, `sibling`

#### Get All Relations
```http
GET /relations
# Filter by house:
GET /relations?house_id=1
```

#### Get Relation by ID
```http
GET /relations/1
```

#### Update Relation
```http
PUT /relations/1
Content-Type: application/json

{
  "relation_type": "spouse"
}
```

#### Delete Relation
```http
DELETE /relations/1
```

### Family Tree

#### Get Family Tree
```http
GET /family-tree/1
```

Returns the complete family tree for a house including all persons and their relationships.

## Database Schema

### Tables

#### admins
- `id` - Primary key
- `username` - Unique username
- `password` - Hashed password
- `created_at` - Timestamp

#### houses
- `id` - Primary key
- `name` - House name
- `created_by` - Foreign key to admins
- `created_at` - Timestamp

#### persons
- `id` - Primary key
- `house_id` - Foreign key to houses
- `name` - Person's name
- `contact` - Contact information
- `description` - Description
- `gender` - 'male' or 'female'
- `dob` - Date of birth
- `created_at` - Timestamp

#### relations
- `id` - Primary key
- `house_id` - Foreign key to houses
- `person_id` - Foreign key to persons
- `related_to_id` - Foreign key to persons
- `relation_type` - 'parent', 'spouse', or 'sibling'
- `created_at` - Timestamp

## Project Structure

```
gofamtree/
├── config/
│   └── db.go              # Database configuration
├── handlers/
│   ├── admin.go           # Admin authentication handlers
│   ├── house.go           # House CRUD handlers
│   ├── person.go          # Person CRUD handlers
│   └── relation.go        # Relation CRUD handlers
├── models/
│   ├── admin.go           # Admin model
│   ├── house.go           # House model
│   ├── person.go          # Person model
│   └── relation.go        # Relation model
├── routes/
│   └── routes.go          # Route definitions
├── utils/
│   └── hash.go            # Password hashing utilities
├── go.mod                 # Go module file
├── main.go                # Application entry point
└── README.md              # This file
```

## Development

### Adding New Features

1. Define the model in `models/`
2. Create handlers in `handlers/`
3. Add routes in `routes/routes.go`
4. Update database migration in `config/db.go`

### Validation Rules

- Duplicate relations are prevented
- Self-relations are not allowed
- Persons must belong to the same house for relations
- Gender must be 'male' or 'female'
- Relation types are restricted to 'parent', 'spouse', 'sibling'

## Environment Variables

- `DATABASE_URL` - PostgreSQL connection string
- `PORT` - Server port (default: 8080)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License. 
# GoFamTree API - Postman Collection Guide

This guide will help you import and use the Postman collection for testing the GoFamTree API.

## 📦 Files Included

- `GoFamTree_API.postman_collection.json` - Main collection with all API endpoints
- `GoFamTree_Local.postman_environment.json` - Environment variables for local development

## 🚀 Quick Setup

### 1. Import Collection and Environment

1. Open Postman
2. Click **Import** button
3. Drag and drop both JSON files:
   - `GoFamTree_API.postman_collection.json`
   - `GoFamTree_Local.postman_environment.json`
4. Select the **GoFamTree - Local** environment from the dropdown

### 2. Start Your API Server

Make sure your GoFamTree API server is running:

```bash
cd ~/workspace/backend/gofamtree
go run main.go
```

The server should start on `http://localhost:8080`

## 📋 Collection Structure

The collection is organized into 5 main folders:

### 🔐 Admin Authentication
- **Register Admin** - Create a new admin account
- **Admin Login** - Login with admin credentials

### 🏠 House Management
- **Create House** - Create a new family house
- **Get All Houses** - List all houses
- **Get House by ID** - Get specific house details
- **Update House** - Modify house information
- **Delete House** - Remove house and all related data

### 👥 Person Management
- **Create Person** - Add new family member
- **Create Person (Female)** - Example for female person
- **Create Person (Child)** - Example for child person
- **Get All Persons** - List all persons
- **Get Persons by House ID** - Filter persons by house
- **Get Person by ID** - Get specific person details
- **Update Person** - Modify person information
- **Delete Person** - Remove person and their relationships

### 🔗 Relationship Management
- **Create Spouse Relation** - Create marriage relationship
- **Create Parent-Child Relation** - Create parent-child bond
- **Create Sibling Relation** - Create sibling relationship
- **Get All Relations** - List all relationships
- **Get Relations by House ID** - Filter relations by house
- **Get Relation by ID** - Get specific relationship
- **Update Relation** - Change relationship type
- **Delete Relation** - Remove relationship

### 🌳 Family Tree
- **Get Family Tree** - Complete family tree for a house

## 🎯 Testing Workflow

Follow this recommended sequence for testing:

### Step 1: Setup Admin
1. Run **Register Admin** to create an admin account
2. Run **Admin Login** to verify authentication

### Step 2: Create House
1. Run **Create House** (use admin ID from Step 1)
2. Note the house ID returned

### Step 3: Add Family Members
1. Run **Create Person** for father
2. Run **Create Person (Female)** for mother
3. Run **Create Person (Child)** for child
4. Note all person IDs returned

### Step 4: Create Relationships
1. Run **Create Spouse Relation** (father + mother)
2. Run **Create Parent-Child Relation** (father + child)
3. Run **Create Parent-Child Relation** (mother + child)

### Step 5: View Family Tree
1. Run **Get Family Tree** to see the complete structure

## 🔧 Environment Variables

The collection uses these environment variables:

| Variable | Default Value | Description |
|----------|---------------|-------------|
| `base_url` | `http://localhost:8080` | API server URL |
| `admin_id` | `1` | Default admin ID |
| `house_id` | `1` | Default house ID |
| `person_id` | `1` | Default person ID |
| `relation_id` | `1` | Default relation ID |

You can modify these in the **GoFamTree - Local** environment as needed.

## 📝 Sample Data Templates

### Admin Registration
```json
{
  "username": "admin",
  "password": "password123"
}
```

### House Creation
```json
{
  "name": "Smith Family",
  "created_by": 1
}
```

### Person Creation
```json
{
  "house_id": 1,
  "name": "John Smith",
  "contact": "john@example.com",
  "description": "Father of the family",
  "gender": "male",
  "dob": "1980-01-15"
}
```

### Relationship Creation
```json
{
  "house_id": 1,
  "person_id": 1,
  "related_to_id": 2,
  "relation_type": "spouse"
}
```

## 🚨 Important Notes

1. **Order Matters**: Create admin → house → persons → relations
2. **IDs**: Update IDs in requests based on actual returned values
3. **Date Format**: Use YYYY-MM-DD for date of birth
4. **Gender**: Only "male" and "female" are accepted
5. **Relation Types**: Only "parent", "spouse", and "sibling" are valid

## 🔍 Troubleshooting

### Common Issues

1. **Connection Refused**: Make sure the API server is running
2. **Admin Not Found**: Register admin first before creating houses
3. **House Not Found**: Create house before adding persons
4. **Invalid Date**: Use YYYY-MM-DD format for dates
5. **Duplicate Relations**: Each person-relation combination must be unique

### Response Codes

- `200` - Success
- `201` - Created successfully
- `400` - Bad request (invalid data)
- `401` - Unauthorized (login required)
- `404` - Not found
- `409` - Conflict (duplicate data)
- `500` - Server error

## 🎉 Happy Testing!

You now have everything needed to test the GoFamTree API. Start with the admin registration and work your way through building a complete family tree!

For additional help, refer to the main README.md file or the API documentation. 
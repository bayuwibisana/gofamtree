#!/bin/bash

# GoFamTree API Test Script
# Make sure the server is running on localhost:8080

BASE_URL="http://localhost:8080"

echo "🌳 GoFamTree API Testing Script"
echo "================================"

# Test 1: Register Admin
echo "📝 1. Registering admin..."
curl -X POST "$BASE_URL/admin/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testadmin",
    "password": "password123"
  }' | jq .
echo ""

# Test 2: Login Admin
echo "🚪 2. Admin login..."
curl -X POST "$BASE_URL/admin/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testadmin",
    "password": "password123"
  }' | jq .
echo ""

# Test 3: Create House
echo "🏠 3. Creating house..."
curl -X POST "$BASE_URL/houses" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smith Family",
    "created_by": 1
  }' | jq .
echo ""

# Test 4: Get All Houses
echo "📋 4. Getting all houses..."
curl -X GET "$BASE_URL/houses" | jq .
echo ""

# Test 5: Create Persons
echo "👥 5. Creating persons..."

# Father
curl -X POST "$BASE_URL/persons" \
  -H "Content-Type: application/json" \
  -d '{
    "house_id": 1,
    "name": "John Smith",
    "contact": "john@example.com",
    "description": "Father of the family",
    "gender": "male",
    "dob": "1980-01-15"
  }' | jq .
echo ""

# Mother
curl -X POST "$BASE_URL/persons" \
  -H "Content-Type: application/json" \
  -d '{
    "house_id": 1,
    "name": "Jane Smith",
    "contact": "jane@example.com",
    "description": "Mother of the family",
    "gender": "female",
    "dob": "1985-03-20"
  }' | jq .
echo ""

# Child
curl -X POST "$BASE_URL/persons" \
  -H "Content-Type: application/json" \
  -d '{
    "house_id": 1,
    "name": "Alice Smith",
    "contact": "alice@example.com",
    "description": "First child",
    "gender": "female",
    "dob": "2010-07-12"
  }' | jq .
echo ""

# Test 6: Get All Persons
echo "📋 6. Getting all persons..."
curl -X GET "$BASE_URL/persons" | jq .
echo ""

# Test 7: Create Relations
echo "🔗 7. Creating relations..."

# Spouse relation
curl -X POST "$BASE_URL/relations" \
  -H "Content-Type: application/json" \
  -d '{
    "house_id": 1,
    "person_id": 1,
    "related_to_id": 2,
    "relation_type": "spouse"
  }' | jq .
echo ""

# Parent-child relation (John -> Alice)
curl -X POST "$BASE_URL/relations" \
  -H "Content-Type: application/json" \
  -d '{
    "house_id": 1,
    "person_id": 1,
    "related_to_id": 3,
    "relation_type": "parent"
  }' | jq .
echo ""

# Parent-child relation (Jane -> Alice)
curl -X POST "$BASE_URL/relations" \
  -H "Content-Type: application/json" \
  -d '{
    "house_id": 1,
    "person_id": 2,
    "related_to_id": 3,
    "relation_type": "parent"
  }' | jq .
echo ""

# Test 8: Get All Relations
echo "📋 8. Getting all relations..."
curl -X GET "$BASE_URL/relations" | jq .
echo ""

# Test 9: Get Family Tree
echo "🌳 9. Getting family tree for house 1..."
curl -X GET "$BASE_URL/family-tree/1" | jq .
echo ""

# Test 10: Update Person
echo "✏️ 10. Updating person..."
curl -X PUT "$BASE_URL/persons/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith Sr.",
    "contact": "johnsr@example.com",
    "description": "Updated father description",
    "gender": "male",
    "dob": "1980-01-15"
  }' | jq .
echo ""

# Test 11: Get Updated Person
echo "👁️ 11. Getting updated person..."
curl -X GET "$BASE_URL/persons/1" | jq .
echo ""

echo "✅ API Testing completed!"
echo "================================" 
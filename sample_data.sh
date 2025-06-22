#!/bin/bash

# Family Tree API Sample Data Script
# This script creates a multi-generational family tree with patriarch to grandchildren

BASE_URL="http://localhost:8080/api/v1"

echo "🌳 Creating Family Tree Sample Data..."
echo "=================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to make API calls and extract ID
create_person() {
    local response=$(curl -s -X POST "$BASE_URL/persons" \
        -H "Content-Type: application/json" \
        -d "$1")
    
    local id=$(echo "$response" | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
    echo "$id"
}

# Function to print person creation
print_creation() {
    echo -e "${GREEN}✓${NC} Created: $1 (ID: $2)"
}

echo -e "${BLUE}Generation 1 - The Patriarchs${NC}"
echo "==============================="

# 1. Create the Patriarch
PATRIARCH_DATA='{
    "first_name": "Robert",
    "last_name": "Johnson",
    "gender": "male",
    "birth_date": "1920-03-15T00:00:00Z",
    "death_date": "2010-11-20T00:00:00Z",
    "biography": "Robert Johnson was a hardworking farmer who built the Johnson family legacy. He served in World War II and returned to establish a successful agricultural business. Known for his wisdom and strong family values.",
    "photo_url": "https://example.com/photos/robert-johnson.jpg"
}'

PATRIARCH_ID=$(create_person "$PATRIARCH_DATA")
print_creation "Robert Johnson (Patriarch)" "$PATRIARCH_ID"

# 2. Create the Matriarch
MATRIARCH_DATA='{
    "first_name": "Eleanor",
    "last_name": "Johnson",
    "gender": "female",
    "birth_date": "1925-07-08T00:00:00Z",
    "death_date": "2015-05-12T00:00:00Z",
    "biography": "Eleanor Johnson was a devoted wife and mother who raised five children while supporting her husband'\''s business. She was known in the community for her charitable work and exceptional cooking.",
    "photo_url": "https://example.com/photos/eleanor-johnson.jpg"
}'

MATRIARCH_ID=$(create_person "$MATRIARCH_DATA")
print_creation "Eleanor Johnson (Matriarch)" "$MATRIARCH_ID"

echo
echo -e "${BLUE}Generation 2 - The Children${NC}"
echo "============================"

# 3. Create First Son - Michael
MICHAEL_DATA='{
    "first_name": "Michael",
    "last_name": "Johnson",
    "gender": "male",
    "birth_date": "1948-01-20T00:00:00Z",
    "biography": "Michael took over the family business and expanded it into real estate. Married to Susan and has three children.",
    "father_id": '$PATRIARCH_ID',
    "mother_id": '$MATRIARCH_ID',
    "photo_url": "https://example.com/photos/michael-johnson.jpg"
}'

MICHAEL_ID=$(create_person "$MICHAEL_DATA")
print_creation "Michael Johnson (Son)" "$MICHAEL_ID"

# 4. Create Michael's Wife - Susan
SUSAN_DATA='{
    "first_name": "Susan",
    "last_name": "Johnson",
    "gender": "female",
    "birth_date": "1950-09-14T00:00:00Z",
    "biography": "Susan is a retired teacher who dedicated her life to education and raising her three children.",
    "photo_url": "https://example.com/photos/susan-johnson.jpg"
}'

SUSAN_ID=$(create_person "$SUSAN_DATA")
print_creation "Susan Johnson (Michael's Wife)" "$SUSAN_ID"

# 5. Create Second Son - David
DAVID_DATA='{
    "first_name": "David",
    "last_name": "Johnson",
    "gender": "male",
    "birth_date": "1950-06-10T00:00:00Z",
    "biography": "David became a doctor and moved to the city. He specializes in cardiology and has two children with his wife Linda.",
    "father_id": '$PATRIARCH_ID',
    "mother_id": '$MATRIARCH_ID',
    "photo_url": "https://example.com/photos/david-johnson.jpg"
}'

DAVID_ID=$(create_person "$DAVID_DATA")
print_creation "David Johnson (Son)" "$DAVID_ID"

# 6. Create David's Wife - Linda
LINDA_DATA='{
    "first_name": "Linda",
    "last_name": "Johnson",
    "gender": "female",
    "birth_date": "1952-11-03T00:00:00Z",
    "biography": "Linda is a successful lawyer who balanced her career with raising two children.",
    "photo_url": "https://example.com/photos/linda-johnson.jpg"
}'

LINDA_ID=$(create_person "$LINDA_DATA")
print_creation "Linda Johnson (David's Wife)" "$LINDA_ID"

# 7. Create Daughter - Sarah
SARAH_DATA='{
    "first_name": "Sarah",
    "last_name": "Williams",
    "gender": "female",
    "birth_date": "1953-04-25T00:00:00Z",
    "biography": "Sarah married Tom Williams and became a nurse. She has four children and is known for her compassionate nature.",
    "father_id": '$PATRIARCH_ID',
    "mother_id": '$MATRIARCH_ID',
    "photo_url": "https://example.com/photos/sarah-williams.jpg"
}'

SARAH_ID=$(create_person "$SARAH_DATA")
print_creation "Sarah Williams (Daughter)" "$SARAH_ID"

# 8. Create Sarah's Husband - Tom
TOM_DATA='{
    "first_name": "Tom",
    "last_name": "Williams",
    "gender": "male",
    "birth_date": "1951-12-18T00:00:00Z",
    "biography": "Tom Williams is a construction contractor who built many homes in the local community.",
    "photo_url": "https://example.com/photos/tom-williams.jpg"
}'

TOM_ID=$(create_person "$TOM_DATA")
print_creation "Tom Williams (Sarah's Husband)" "$TOM_ID"

echo
echo -e "${BLUE}Generation 3 - The Grandchildren${NC}"
echo "================================="

# Michael & Susan's Children
JAMES_DATA='{
    "first_name": "James",
    "last_name": "Johnson",
    "gender": "male",
    "birth_date": "1975-03-12T00:00:00Z",
    "biography": "James is a software engineer working in tech. He enjoys hiking and photography.",
    "father_id": '$MICHAEL_ID',
    "mother_id": '$SUSAN_ID',
    "photo_url": "https://example.com/photos/james-johnson.jpg"
}'

JAMES_ID=$(create_person "$JAMES_DATA")
print_creation "James Johnson (Michael's Son)" "$JAMES_ID"

EMILY_DATA='{
    "first_name": "Emily",
    "last_name": "Johnson",
    "gender": "female",
    "birth_date": "1977-08-05T00:00:00Z",
    "biography": "Emily is a marketing manager who loves travel and cooking. She'\''s planning to start her own consulting business.",
    "father_id": '$MICHAEL_ID',
    "mother_id": '$SUSAN_ID',
    "photo_url": "https://example.com/photos/emily-johnson.jpg"
}'

EMILY_ID=$(create_person "$EMILY_DATA")
print_creation "Emily Johnson (Michael's Daughter)" "$EMILY_ID"

CHRISTOPHER_DATA='{
    "first_name": "Christopher",
    "last_name": "Johnson",
    "gender": "male",
    "birth_date": "1980-01-22T00:00:00Z",
    "biography": "Christopher is an artist and graphic designer. He has exhibited his work in several galleries.",
    "father_id": '$MICHAEL_ID',
    "mother_id": '$SUSAN_ID',
    "photo_url": "https://example.com/photos/christopher-johnson.jpg"
}'

CHRISTOPHER_ID=$(create_person "$CHRISTOPHER_DATA")
print_creation "Christopher Johnson (Michael's Son)" "$CHRISTOPHER_ID"

# David & Linda's Children
AMANDA_DATA='{
    "first_name": "Amanda",
    "last_name": "Johnson",
    "gender": "female",
    "birth_date": "1978-05-30T00:00:00Z",
    "biography": "Amanda followed in her father'\''s footsteps and became a pediatrician. She works at a children'\''s hospital.",
    "father_id": '$DAVID_ID',
    "mother_id": '$LINDA_ID',
    "photo_url": "https://example.com/photos/amanda-johnson.jpg"
}'

AMANDA_ID=$(create_person "$AMANDA_DATA")
print_creation "Amanda Johnson (David's Daughter)" "$AMANDA_ID"

RYAN_DATA='{
    "first_name": "Ryan",
    "last_name": "Johnson",
    "gender": "male",
    "birth_date": "1981-11-15T00:00:00Z",
    "biography": "Ryan is a financial advisor who helps families plan for their future. He'\''s passionate about financial literacy.",
    "father_id": '$DAVID_ID',
    "mother_id": '$LINDA_ID',
    "photo_url": "https://example.com/photos/ryan-johnson.jpg"
}'

RYAN_ID=$(create_person "$RYAN_DATA")
print_creation "Ryan Johnson (David's Son)" "$RYAN_ID"

# Sarah & Tom's Children
MELISSA_DATA='{
    "first_name": "Melissa",
    "last_name": "Williams",
    "gender": "female",
    "birth_date": "1976-07-18T00:00:00Z",
    "biography": "Melissa is a veterinarian who runs her own animal clinic. She has a special love for rescue animals.",
    "father_id": '$TOM_ID',
    "mother_id": '$SARAH_ID',
    "photo_url": "https://example.com/photos/melissa-williams.jpg"
}'

MELISSA_ID=$(create_person "$MELISSA_DATA")
print_creation "Melissa Williams (Sarah's Daughter)" "$MELISSA_ID"

DANIEL_DATA='{
    "first_name": "Daniel",
    "last_name": "Williams",
    "gender": "male",
    "birth_date": "1979-02-28T00:00:00Z",
    "biography": "Daniel is a high school history teacher and basketball coach. He'\''s beloved by his students.",
    "father_id": '$TOM_ID',
    "mother_id": '$SARAH_ID',
    "photo_url": "https://example.com/photos/daniel-williams.jpg"
}'

DANIEL_ID=$(create_person "$DANIEL_DATA")
print_creation "Daniel Williams (Sarah's Son)" "$DANIEL_ID"

NICOLE_DATA='{
    "first_name": "Nicole",
    "last_name": "Williams",
    "gender": "female",
    "birth_date": "1982-10-12T00:00:00Z",
    "biography": "Nicole is a social worker dedicated to helping at-risk youth. She'\''s working on her master'\''s degree.",
    "father_id": '$TOM_ID',
    "mother_id": '$SARAH_ID',
    "photo_url": "https://example.com/photos/nicole-williams.jpg"
}'

NICOLE_ID=$(create_person "$NICOLE_DATA")
print_creation "Nicole Williams (Sarah's Daughter)" "$NICOLE_ID"

BRANDON_DATA='{
    "first_name": "Brandon",
    "last_name": "Williams",
    "gender": "male",
    "birth_date": "1985-12-05T00:00:00Z",
    "biography": "Brandon is a chef who owns a small restaurant. He specializes in farm-to-table cuisine.",
    "father_id": '$TOM_ID',
    "mother_id": '$SARAH_ID',
    "photo_url": "https://example.com/photos/brandon-williams.jpg"
}'

BRANDON_ID=$(create_person "$BRANDON_DATA")
print_creation "Brandon Williams (Sarah's Son)" "$BRANDON_ID"

echo
echo -e "${YELLOW}📊 Family Tree Summary${NC}"
echo "======================"
echo -e "${GREEN}✓${NC} Generation 1: 2 people (Robert & Eleanor Johnson)"
echo -e "${GREEN}✓${NC} Generation 2: 5 people (3 children + 2 spouses)"
echo -e "${GREEN}✓${NC} Generation 3: 9 grandchildren"
echo -e "${GREEN}✓${NC} Total: 16 people created"
echo
echo -e "${BLUE}🔍 Test the Family Tree:${NC}"
echo "curl $BASE_URL/persons/$PATRIARCH_ID/family-tree"
echo
echo -e "${BLUE}🔍 Search for family members:${NC}"
echo "curl \"$BASE_URL/persons/search?q=Johnson\""
echo
echo -e "${GREEN}✅ Sample data creation completed!${NC}" 
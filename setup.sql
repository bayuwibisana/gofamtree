-- GoFamTree Database Setup Script
-- Run this in your PostgreSQL instance

-- Create database
CREATE DATABASE gofamtree_new;

-- Connect to the database
\c gofamtree_new;

-- Create tables (these will be automatically created by GORM, but here for reference)

-- Admins table
CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL, -- hashed
    created_at TIMESTAMP DEFAULT NOW()
);

-- Houses table
CREATE TABLE IF NOT EXISTS houses (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_by INTEGER NOT NULL REFERENCES admins(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Persons table
CREATE TABLE IF NOT EXISTS persons (
    id SERIAL PRIMARY KEY,
    house_id INTEGER NOT NULL REFERENCES houses(id),
    name TEXT NOT NULL,
    contact TEXT,
    description TEXT,
    gender TEXT CHECK (gender IN ('male', 'female')),
    dob DATE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Relations table
CREATE TABLE IF NOT EXISTS relations (
    id SERIAL PRIMARY KEY,
    house_id INTEGER NOT NULL REFERENCES houses(id),
    person_id INTEGER NOT NULL REFERENCES persons(id),
    related_to_id INTEGER NOT NULL REFERENCES persons(id),
    relation_type TEXT CHECK (relation_type IN ('parent', 'spouse', 'sibling')),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(person_id, related_to_id, relation_type)
);

-- Indexes for better performance
CREATE INDEX idx_houses_created_by ON houses(created_by);
CREATE INDEX idx_persons_house_id ON persons(house_id);
CREATE INDEX idx_relations_house_id ON relations(house_id);
CREATE INDEX idx_relations_person_id ON relations(person_id);
CREATE INDEX idx_relations_related_to_id ON relations(related_to_id);

-- Sample data (optional)
-- INSERT INTO admins (username, password) VALUES ('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'); -- password: 'password'

PRINT 'Database setup completed successfully!'; 
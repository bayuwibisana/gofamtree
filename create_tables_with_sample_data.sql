-- GoFamTree Database Setup with 4 Generations Sample Data
-- Connect to gofamtree_new database before running this script

-- Drop tables if they exist (for clean setup)
DROP TABLE IF EXISTS relations CASCADE;
DROP TABLE IF EXISTS persons CASCADE;
DROP TABLE IF EXISTS houses CASCADE;
DROP TABLE IF EXISTS admins CASCADE;

-- Create tables
CREATE TABLE admins (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL, -- hashed
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE houses (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_by INTEGER NOT NULL REFERENCES admins(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE persons (
    id SERIAL PRIMARY KEY,
    house_id INTEGER NOT NULL REFERENCES houses(id),
    name TEXT NOT NULL,
    contact TEXT,
    description TEXT,
    gender TEXT CHECK (gender IN ('male', 'female')),
    dob DATE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE relations (
    id SERIAL PRIMARY KEY,
    house_id INTEGER NOT NULL REFERENCES houses(id),
    person_id INTEGER NOT NULL REFERENCES persons(id),
    related_to_id INTEGER NOT NULL REFERENCES persons(id),
    relation_type TEXT CHECK (relation_type IN ('parent', 'spouse', 'sibling')),
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(person_id, related_to_id, relation_type)
);

-- Create indexes for better performance
CREATE INDEX idx_houses_created_by ON houses(created_by);
CREATE INDEX idx_persons_house_id ON persons(house_id);
CREATE INDEX idx_relations_house_id ON relations(house_id);
CREATE INDEX idx_relations_person_id ON relations(person_id);
CREATE INDEX idx_relations_related_to_id ON relations(related_to_id);

-- Insert sample admin (password is 'password123' hashed with bcrypt)
INSERT INTO admins (username, password, created_at) VALUES 
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW());

-- Insert sample house
INSERT INTO houses (name, created_by, created_at) VALUES 
('Johnson Family Dynasty', 1, NOW());

-- Insert 4 Generations of Sample Data
-- Generation 1: Great-Grandparents (Born 1920s)
INSERT INTO persons (house_id, name, contact, description, gender, dob, created_at) VALUES 
(1, 'William Johnson Sr.', 'william.sr@example.com', 'Great-grandfather, family patriarch', 'male', '1925-03-15', NOW()),
(1, 'Mary Johnson', 'mary.johnson@example.com', 'Great-grandmother, beloved matriarch', 'female', '1928-07-22', NOW());

-- Generation 2: Grandparents (Born 1940s-1950s)
INSERT INTO persons (house_id, name, contact, description, gender, dob, created_at) VALUES 
(1, 'Robert Johnson', 'robert.johnson@example.com', 'Eldest son, engineer', 'male', '1948-01-10', NOW()),
(1, 'Linda Johnson', 'linda.johnson@example.com', 'Roberts wife, teacher', 'female', '1952-09-18', NOW()),
(1, 'James Johnson', 'james.johnson@example.com', 'Younger son, doctor', 'male', '1950-11-05', NOW()),
(1, 'Patricia Johnson', 'patricia.johnson@example.com', 'James wife, nurse', 'female', '1954-04-12', NOW());

-- Generation 3: Parents (Born 1970s-1980s)
INSERT INTO persons (house_id, name, contact, description, gender, dob, created_at) VALUES 
(1, 'Michael Johnson', 'michael.johnson@example.com', 'Roberts son, software developer', 'male', '1975-06-20', NOW()),
(1, 'Sarah Johnson', 'sarah.johnson@example.com', 'Michaels wife, marketing manager', 'female', '1978-12-03', NOW()),
(1, 'David Johnson', 'david.johnson@example.com', 'James son, lawyer', 'male', '1973-08-14', NOW()),
(1, 'Jennifer Johnson', 'jennifer.johnson@example.com', 'Davids wife, graphic designer', 'female', '1976-02-28', NOW());

-- Generation 4: Children (Born 2000s-2010s)
INSERT INTO persons (house_id, name, contact, description, gender, dob, created_at) VALUES 
(1, 'Christopher Johnson', 'chris.johnson@example.com', 'Michael and Sarahs eldest son', 'male', '2005-04-15', NOW()),
(1, 'Emily Johnson', 'emily.johnson@example.com', 'Michael and Sarahs daughter', 'female', '2008-09-22', NOW()),
(1, 'Matthew Johnson', 'matthew.johnson@example.com', 'Michael and Sarahs youngest son', 'male', '2012-01-08', NOW()),
(1, 'Olivia Johnson', 'olivia.johnson@example.com', 'David and Jennifers daughter', 'female', '2003-11-30', NOW()),
(1, 'Daniel Johnson', 'daniel.johnson@example.com', 'David and Jennifers eldest son', 'male', '2006-07-18', NOW()),
(1, 'Sophia Johnson', 'sophia.johnson@example.com', 'David and Jennifers youngest daughter', 'female', '2010-03-12', NOW());

-- Create Relationships
-- Generation 1: Great-Grandparents (Spouse relationship)
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 1, 2, 'spouse', NOW()),  -- William Sr. ↔ Mary
(1, 2, 1, 'spouse', NOW());

-- Generation 1 → Generation 2 (Parent-Child relationships)
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 1, 3, 'parent', NOW()),  -- William Sr. → Robert
(1, 2, 3, 'parent', NOW()),  -- Mary → Robert
(1, 1, 5, 'parent', NOW()),  -- William Sr. → James
(1, 2, 5, 'parent', NOW());  -- Mary → James

-- Generation 2: Grandparents (Spouse relationships)
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 3, 4, 'spouse', NOW()),  -- Robert ↔ Linda
(1, 4, 3, 'spouse', NOW()),
(1, 5, 6, 'spouse', NOW()),  -- James ↔ Patricia
(1, 6, 5, 'spouse', NOW());

-- Generation 2: Siblings
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 3, 5, 'sibling', NOW()),  -- Robert ↔ James
(1, 5, 3, 'sibling', NOW());

-- Generation 2 → Generation 3 (Parent-Child relationships)
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 3, 7, 'parent', NOW()),  -- Robert → Michael
(1, 4, 7, 'parent', NOW()),  -- Linda → Michael
(1, 5, 9, 'parent', NOW()),  -- James → David
(1, 6, 9, 'parent', NOW());  -- Patricia → David

-- Generation 3: Parents (Spouse relationships)
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 7, 8, 'spouse', NOW()),  -- Michael ↔ Sarah
(1, 8, 7, 'spouse', NOW()),
(1, 9, 10, 'spouse', NOW()), -- David ↔ Jennifer
(1, 10, 9, 'spouse', NOW());

-- Generation 3: Cousins (represented as siblings of parents)
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 7, 9, 'sibling', NOW()),  -- Michael ↔ David (cousins, but using sibling for simplicity)
(1, 9, 7, 'sibling', NOW());

-- Generation 3 → Generation 4 (Parent-Child relationships)
-- Michael & Sarah's children
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 7, 11, 'parent', NOW()), -- Michael → Christopher
(1, 8, 11, 'parent', NOW()), -- Sarah → Christopher
(1, 7, 12, 'parent', NOW()), -- Michael → Emily
(1, 8, 12, 'parent', NOW()), -- Sarah → Emily
(1, 7, 13, 'parent', NOW()), -- Michael → Matthew
(1, 8, 13, 'parent', NOW()); -- Sarah → Matthew

-- David & Jennifer's children
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 9, 14, 'parent', NOW()),  -- David → Olivia
(1, 10, 14, 'parent', NOW()), -- Jennifer → Olivia
(1, 9, 15, 'parent', NOW()),  -- David → Daniel
(1, 10, 15, 'parent', NOW()), -- Jennifer → Daniel
(1, 9, 16, 'parent', NOW()),  -- David → Sophia
(1, 10, 16, 'parent', NOW()); -- Jennifer → Sophia

-- Generation 4: Siblings within families
-- Michael & Sarah's children siblings
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 11, 12, 'sibling', NOW()), -- Christopher ↔ Emily
(1, 12, 11, 'sibling', NOW()),
(1, 11, 13, 'sibling', NOW()), -- Christopher ↔ Matthew
(1, 13, 11, 'sibling', NOW()),
(1, 12, 13, 'sibling', NOW()), -- Emily ↔ Matthew
(1, 13, 12, 'sibling', NOW());

-- David & Jennifer's children siblings
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 14, 15, 'sibling', NOW()), -- Olivia ↔ Daniel
(1, 15, 14, 'sibling', NOW()),
(1, 14, 16, 'sibling', NOW()), -- Olivia ↔ Sophia
(1, 16, 14, 'sibling', NOW()),
(1, 15, 16, 'sibling', NOW()), -- Daniel ↔ Sophia
(1, 16, 15, 'sibling', NOW());

-- Generation 4: Cousins (siblings relationship for simplicity)
INSERT INTO relations (house_id, person_id, related_to_id, relation_type, created_at) VALUES 
(1, 11, 14, 'sibling', NOW()), -- Christopher ↔ Olivia (cousins)
(1, 14, 11, 'sibling', NOW()),
(1, 12, 15, 'sibling', NOW()), -- Emily ↔ Daniel (cousins)
(1, 15, 12, 'sibling', NOW()),
(1, 13, 16, 'sibling', NOW()), -- Matthew ↔ Sophia (cousins)
(1, 16, 13, 'sibling', NOW());

-- Display summary
DO $$
BEGIN
    RAISE NOTICE 'Database setup completed successfully!';
    RAISE NOTICE 'Created:';
    RAISE NOTICE '- % admins', (SELECT COUNT(*) FROM admins);
    RAISE NOTICE '- % houses', (SELECT COUNT(*) FROM houses);
    RAISE NOTICE '- % persons (4 generations)', (SELECT COUNT(*) FROM persons);
    RAISE NOTICE '- % relations', (SELECT COUNT(*) FROM relations);
    RAISE NOTICE '';
    RAISE NOTICE 'Family Structure:';
    RAISE NOTICE 'Generation 1: William Sr. & Mary (Great-grandparents)';
    RAISE NOTICE 'Generation 2: Robert & Linda, James & Patricia (Grandparents)';
    RAISE NOTICE 'Generation 3: Michael & Sarah, David & Jennifer (Parents)';
    RAISE NOTICE 'Generation 4: 6 children/grandchildren';
END $$; 
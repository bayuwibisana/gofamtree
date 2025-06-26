package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() {
	// Use environment variable for DSN or default without password
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres dbname=gofamtree_new port=5432 sslmode=disable"
	}
	
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Since we're managing the schema manually with SQL scripts,
	// we'll disable auto-migration to avoid conflicts
	// 
	// If you want to enable auto-migration, comment out the lines below
	// and uncomment the AutoMigrate section

	log.Println("Using manual schema management - auto-migration disabled")
	
	// Test the connection by checking if tables exist
	var tableExists bool
	result := DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'admins')").Scan(&tableExists)
	if result.Error != nil {
		log.Fatal("Failed to check database schema:", result.Error)
	}
	
	if !tableExists {
		log.Println("WARNING: Database tables not found. Please run the SQL setup script:")
		log.Println("psql -d gofamtree_new -f create_tables_with_sample_data.sql")
	} else {
		log.Println("Database schema verified - tables exist")
	}

	// Uncomment this section if you want to use GORM auto-migration instead of manual SQL:
	/*
	err = DB.AutoMigrate(
		&models.Admin{},
		&models.House{},
		&models.Person{},
		&models.Relation{},
	)
	
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	
	log.Println("Database migration completed")
	*/
}

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/bayuwibisana/gofamtree/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Set defaults if environment variables are empty
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if dbname == "" {
		dbname = "gofamtree"
		log.Printf("Warning: DB_NAME was empty, using default: gofamtree")
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	// Use PostgreSQL URL format for reliable connection
	var dsn string
	if password != "" {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=Asia/Jakarta",
			user, password, host, port, dbname, sslmode)
	} else {
		dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s&TimeZone=Asia/Jakarta",
			user, host, port, dbname, sslmode)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Verify which database we're connected to
	var currentDB string
	err = DB.Raw("SELECT current_database()").Scan(&currentDB).Error
	if err != nil {
		log.Fatal("Failed to verify database connection:", err)
	}
	log.Printf("✅ Connected to database: %s", currentDB)

	// Auto-migrate the schema
	err = DB.AutoMigrate(&models.Person{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")
}

func GetDB() *gorm.DB {
	return DB
} 
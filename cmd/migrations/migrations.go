package main

import (
	"absen-backend/config"
	"absen-backend/migrations"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Usage: go run cmd/migrations/migrations.go
	config.ConnectDB()

	// ambil dari global var
	db := config.DB

	migs := []func(*mongo.Database) error{
		migrations.Users,
		migrations.Roleusers,
	}

	fmt.Println("🚀 Running migrations...")
	for _, mig := range migs {
		start := time.Now()
		if err := mig(db); err != nil {
			log.Fatalf("❌ Migration failed: %v\n", err)
		}
		fmt.Printf("✅ Migration %T completed in %v\n", mig, time.Since(start))
	}
	fmt.Println("🎉 All migrations completed successfully!")
}

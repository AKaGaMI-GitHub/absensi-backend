package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/migrate.go <migration_name>")
		return
	}

	// ambil nama migration dari argument
	name := os.Args[1]
	timestamp := time.Now().Format("20060102150405") // format YYYYMMDDHHMMSS
	fileName := fmt.Sprintf("migrations/%s_%s.go", timestamp, name)

	// isi template migration
	content := fmt.Sprintf(`package migrations

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func %s(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collectionName := "<migration_name>"
	// TODO: implement migration logic here (create collection, seed, index, etc.)

	return nil
}
`, toCamelCase(name))

	// buat file
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error creating migration file:", err)
		return
	}

	fmt.Println("Migration created:", fileName)
}

// ubah migration_name jadi CamelCase
func toCamelCase(input string) string {
	words := strings.Split(input, "_")
	for i := range words {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}

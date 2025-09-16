package migrations

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Users(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// TODO: implement migration logic here (create collection, seed, index, etc.)
	// Nama collection
	collectionName := "users"

	// Cek apakah collection sudah ada
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	// Kalau belum ada, buat collection
	if len(collections) == 0 {
		if err := db.CreateCollection(ctx, collectionName); err != nil {
			return fmt.Errorf("failed to create collection %s: %w", collectionName, err)
		}
	}

	// Tambahkan index unik untuk email
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}}, // ascending
		Options: options.Index().SetUnique(true).SetName("unique_email"),
	}

	_, err = db.Collection(collectionName).Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	return nil
}

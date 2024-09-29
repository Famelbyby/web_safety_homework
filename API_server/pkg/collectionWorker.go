package pkg

import (
	"context"
	"main/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCurrentHTTPRecordByID(collection *mongo.Collection, id int) (domain.HTTPEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var result domain.HTTPEntity
	err := collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)

	return result, err
}

func GetCurrentHTTPSRecordByID(collection *mongo.Collection, id int) (domain.HTTPSEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var result domain.HTTPSEntity
	err := collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&result)

	return result, err
}

func GetAllHTTPRecords(collection *mongo.Collection) ([]domain.HTTPEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	var results []domain.HTTPEntity

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func GetAllHTTPSRecords(collection *mongo.Collection) ([]domain.HTTPSEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	var results []domain.HTTPSEntity

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

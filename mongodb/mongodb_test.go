package mongodb

import (
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestNewDB(t *testing.T) {
	_, err := NewDB("mongodb://admin:password@localhost:27017", "test", 10*time.Second)
	if err != nil {
		panic(err)
	}
}

func TestInsertOne(t *testing.T) {
	db, err := NewDB("mongodb://admin:password@localhost:27017", "test", 10*time.Second)
	if err != nil {
		panic(err)
	}

	var docs []any
	for i := 0; i < 10; i++ {
		doc := bson.D{{"name", "test"}, {"age", i}}
		docs = append(docs, doc)
	}

	err = db.BatchInsert("test", docs)
	if err != nil {
		log.Fatalf("Failed to insert documents: %v", err)
		return
	}
}

func TestReplace(t *testing.T) {
	db, err := NewDB("mongodb://admin:password@localhost:27017", "test", 10*time.Second)
	if err != nil {
		panic(err)
	}

	filter := bson.D{{"name", "test"}}
	replacement := bson.D{{"name", "Bob"}, {"age", 30}}

	err = db.Replace("test", filter, replacement)
	if err != nil {
		log.Fatalf("Failed to replace document: %v", err)
		return
	}
}

func TestUpdate(t *testing.T) {
	db, err := NewDB("mongodb://admin:password@localhost:27017", "test", 10*time.Second)
	if err != nil {
		panic(err)
	}

	filter := bson.D{{"name", "test"}}
	// 更新
	update := bson.D{{"$set", bson.D{{"name", "Bob"}}}}
	// 增加
	// update := bson.D{{"$inc", bson.D{{"age", 30}}}}

	err = db.Update("test", filter, update)
	if err != nil {
		log.Fatalf("Failed to replace document: %v", err)
		return
	}
}

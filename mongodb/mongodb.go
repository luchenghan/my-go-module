package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	BatchInsert(collection string, docs []any) error
	Replace(collection string, filter any, replacement any) error
	Update(collection string, filter any, update interface{}) error
	Close() error
}

type mongodb struct {
	timeout time.Duration
	*mongo.Database
}

func NewDB(uri string, dbName string, timeout time.Duration) (DB, error) {
	db := new(mongodb)
	db.timeout = timeout

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db.Database = client.Database(dbName)

	return db, nil
}

func (m *mongodb) BatchInsert(collection string, docs []any) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.Database.Collection(collection).InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongodb) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	err := m.Database.Client().Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongodb) Replace(collection string, filter any, replacement any) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.Database.Collection(collection).ReplaceOne(ctx, filter, replacement)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongodb) Update(collection string, filter any, update any) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	_, err := m.Database.Collection(collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

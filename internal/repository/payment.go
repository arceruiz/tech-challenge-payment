package repository

import (
	"context"
	"errors"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collection = "payment"
	database   = "payment"
)

var (
	cfg           = &config.Cfg
	ErrorNotFound = errors.New("entity not found")
)

func NewMongo() *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.DB.ConnectionString))
	if err != nil {
		panic(err)
	}
	db := client.Database(database)
	return db
}

type PaymentRepository interface {
	GetByID(context.Context, string) (*canonical.Payment, error)
	Update(ctx context.Context, id string, payment canonical.Payment) error
	Create(ctx context.Context, payment canonical.Payment) (canonical.Payment, error)
	GetAll(ctx context.Context) ([]canonical.Payment, error)
}

type paymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepo(db *mongo.Database) PaymentRepository {
	return &paymentRepository{
		collection: db.Collection(collection),
	}
}

func (r *paymentRepository) Create(ctx context.Context, payment canonical.Payment) (canonical.Payment, error) {

	_, err := r.collection.InsertOne(ctx, payment)
	if err != nil {
		return payment, err
	}
	return payment, nil

}

func (r *paymentRepository) Update(ctx context.Context, id string, payment canonical.Payment) error {
	filter := bson.M{"_id": id}
	fields := bson.M{"$set": payment}

	_, err := r.collection.UpdateOne(ctx, filter, fields)
	if err != nil {
		return err
	}
	return nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id string) (*canonical.Payment, error) {

	var payment canonical.Payment

	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&payment)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) GetAll(ctx context.Context) ([]canonical.Payment, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	var results []canonical.Payment
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

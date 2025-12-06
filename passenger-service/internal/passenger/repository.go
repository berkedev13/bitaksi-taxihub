package passenger

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Insert(ctx context.Context, p *Passenger) (primitive.ObjectID, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Passenger, error)
	Update(ctx context.Context, p *Passenger) error
	List(ctx context.Context, page, pageSize int64) ([]Passenger, error)
	FindAll(ctx context.Context) ([]Passenger, error)
}

type repository struct {
	collection *mongo.Collection
}

func NewRepository(col *mongo.Collection) Repository {
	return &repository{collection: col}
}

func (r *repository) Insert(ctx context.Context, p *Passenger) (primitive.ObjectID, error) {
	res, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return primitive.NilObjectID, err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, mongo.ErrClientDisconnected
	}
	return oid, nil
}

func (r *repository) FindByID(ctx context.Context, id primitive.ObjectID) (*Passenger, error) {
	var p Passenger
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) Update(ctx context.Context, p *Passenger) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": p.ID}, p)
	return err
}

func (r *repository) List(ctx context.Context, page, pageSize int64) ([]Passenger, error) {
	skip := (page - 1) * pageSize

	opts := options.Find()
	opts.Skip = &skip
	opts.Limit = &pageSize

	cur, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var passengers []Passenger
	if err := cur.All(ctx, &passengers); err != nil {
		return nil, err
	}
	return passengers, nil
}

func (r *repository) FindAll(ctx context.Context) ([]Passenger, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var passengers []Passenger
	if err := cur.All(ctx, &passengers); err != nil {
		return nil, err
	}
	return passengers, nil
}

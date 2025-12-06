package driver

import (
	"context"
	"errors"
	"time"

	"github.com/berkedev13/bitaksi-driver-service/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	InsertDriver(ctx context.Context, d *Driver) (primitive.ObjectID, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*Driver, error)
	UpdateDriver(ctx context.Context, d *Driver) error
	ListDrivers(ctx context.Context, page, pageSize int64) ([]Driver, error)
	FindByTaxiType(ctx context.Context, taxiType string) ([]Driver, error)
}

type mongoRepository struct {
	conn *db.MongoConnection
}

func NewRepository(conn *db.MongoConnection) Repository {
	return &mongoRepository{conn: conn}
}

func (r *mongoRepository) InsertDriver(ctx context.Context, d *Driver) (primitive.ObjectID, error) {
	d.CreatedAt = time.Now().UTC()
	d.UpdatedAt = time.Now().UTC()

	res, err := r.conn.DriverColl.InsertOne(ctx, d)
	if err != nil {
		return primitive.NilObjectID, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("insertedID is not an ObjectID")
	}

	return oid, nil
}

func (r *mongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Driver, error) {
	var d Driver
	err := r.conn.DriverColl.FindOne(ctx, bson.M{"_id": id}).Decode(&d)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *mongoRepository) UpdateDriver(ctx context.Context, d *Driver) error {
	d.UpdatedAt = time.Now().UTC()

	_, err := r.conn.DriverColl.ReplaceOne(ctx, bson.M{"_id": d.ID}, d)
	return err
}

func (r *mongoRepository) ListDrivers(ctx context.Context, page, pageSize int64) ([]Driver, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	skip := (page - 1) * pageSize

	opts := options.Find().
		SetSkip(skip).
		SetLimit(pageSize).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := r.conn.DriverColl.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var drivers []Driver
	for cur.Next(ctx) {
		var d Driver
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return drivers, nil
}

func (r *mongoRepository) FindByTaxiType(ctx context.Context, taxiType string) ([]Driver, error) {
	filter := bson.M{}
	if taxiType != "" {
		filter["taxiType"] = taxiType
	}

	cur, err := r.conn.DriverColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var drivers []Driver
	for cur.Next(ctx) {
		var d Driver
		if err := cur.Decode(&d); err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return drivers, nil
}

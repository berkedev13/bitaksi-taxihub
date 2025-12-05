package driver

import (
	"context"
	"errors"
	"time"

	"github.com/berkedev13/bitaksi-driver-service/internal/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	InsertDriver(ctx context.Context, d *Driver) (primitive.ObjectID, error)
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

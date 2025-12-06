package passenger

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Passenger struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
	Phone     string             `bson:"phone" json:"phone"`
	Location  Location           `bson:"location" json:"location"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Location struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lon float64 `bson:"lon" json:"lon"`
}

type CreatePassengerRequest struct {
	FirstName string  `json:"firstName" binding:"required"`
	LastName  string  `json:"lastName" binding:"required"`
	Phone     string  `json:"phone" binding:"required"`
	Lat       float64 `json:"lat" binding:"required"`
	Lon       float64 `json:"lon" binding:"required"`
}

type UpdatePassengerRequest struct {
	FirstName *string  `json:"firstName"`
	LastName  *string  `json:"lastName"`
	Phone     *string  `json:"phone"`
	Lat       *float64 `json:"lat"`
	Lon       *float64 `json:"lon"`
}

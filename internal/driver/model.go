package driver

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lon float64 `json:"lon" bson:"lon"`
}

type Driver struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Plate     string             `json:"plate" bson:"plate"`
	TaxiType  string             `json:"taxiType" bson:"taxiType"`
	CarBrand  string             `json:"carBrand" bson:"carBrand"`
	CarModel  string             `json:"carModel" bson:"carModel"`
	Location  Location           `json:"location" bson:"location"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type CreateDriverRequest struct {
	FirstName string  `json:"firstName" binding:"required"`
	LastName  string  `json:"lastName" binding:"required"`
	Plate     string  `json:"plate" binding:"required"`
	TaxiType  string  `json:"taxiType" binding:"required"`
	CarBrand  string  `json:"carBrand" binding:"required"`
	CarModel  string  `json:"carModel" binding:"required"`
	Lat       float64 `json:"lat" binding:"required"`
	Lon       float64 `json:"lon" binding:"required"`
}

package passenger

import (
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrPassengerNotFound = errors.New("passenger not found")

type Service interface {
	Create(ctx context.Context, req CreatePassengerRequest) (*Passenger, error)
	Update(ctx context.Context, id string, req UpdatePassengerRequest) (*Passenger, error)
	List(ctx context.Context, page, pageSize int64) ([]Passenger, error)
	GetNearby(ctx context.Context, lat, lon float64) ([]Passenger, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreatePassengerRequest) (*Passenger, error) {
	now := time.Now().UTC()
	p := &Passenger{
		ID:        primitive.NewObjectID(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Location: Location{
			Lat: req.Lat,
			Lon: req.Lon,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err := s.repo.Insert(ctx, p)
	return p, err
}

func (s *service) Update(ctx context.Context, id string, req UpdatePassengerRequest) (*Passenger, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrPassengerNotFound
	}

	if req.FirstName != nil {
		existing.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		existing.LastName = *req.LastName
	}
	if req.Phone != nil {
		existing.Phone = *req.Phone
	}
	if req.Lat != nil {
		existing.Location.Lat = *req.Lat
	}
	if req.Lon != nil {
		existing.Location.Lon = *req.Lon
	}

	existing.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *service) List(ctx context.Context, page, pageSize int64) ([]Passenger, error) {
	return s.repo.List(ctx, page, pageSize)
}

func (s *service) GetNearby(ctx context.Context, lat, lon float64) ([]Passenger, error) {
	all, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []Passenger
	for _, p := range all {
		dist := haversineKm(lat, lon, p.Location.Lat, p.Location.Lon)
		if dist <= 6.0 {
			result = append(result, p)
		}
	}
	return result, nil
}

func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0
	rad := math.Pi / 180.0

	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad
	lat1R := lat1 * rad
	lat2R := lat2 * rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1R)*math.Cos(lat2R)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

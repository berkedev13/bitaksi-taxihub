package driver

import (
	"context"
	"errors"
	"math"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrDriverNotFound = errors.New("driver not found")

type Service interface {
	CreateDriver(ctx context.Context, req CreateDriverRequest) (*Driver, error)
	UpdateDriver(ctx context.Context, id string, req UpdateDriverRequest) (*Driver, error)
	ListDrivers(ctx context.Context, page, pageSize int64) ([]Driver, error)
	GetNearbyDrivers(ctx context.Context, lat, lon float64, taxiType string) ([]NearbyDriver, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateDriver(ctx context.Context, req CreateDriverRequest) (*Driver, error) {
	d := &Driver{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Plate:     req.Plate,
		TaxiType:  req.TaxiType,
		CarBrand:  req.CarBrand,
		CarModel:  req.CarModel,
		Location: Location{
			Lat: req.Lat,
			Lon: req.Lon,
		},
	}

	id, err := s.repo.InsertDriver(ctx, d)
	if err != nil {
		return nil, err
	}

	d.ID = id
	return d, nil
}

func (s *service) UpdateDriver(ctx context.Context, id string, req UpdateDriverRequest) (*Driver, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}

	existing, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrDriverNotFound
	}

	if req.FirstName != nil {
		existing.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		existing.LastName = *req.LastName
	}
	if req.Plate != nil {
		existing.Plate = *req.Plate
	}
	if req.TaxiType != nil {
		existing.TaxiType = *req.TaxiType
	}
	if req.CarBrand != nil {
		existing.CarBrand = *req.CarBrand
	}
	if req.CarModel != nil {
		existing.CarModel = *req.CarModel
	}
	if req.Lat != nil {
		existing.Location.Lat = *req.Lat
	}
	if req.Lon != nil {
		existing.Location.Lon = *req.Lon
	}

	existing.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateDriver(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *service) ListDrivers(ctx context.Context, page, pageSize int64) ([]Driver, error) {
	return s.repo.ListDrivers(ctx, page, pageSize)
}

func (s *service) GetNearbyDrivers(ctx context.Context, lat, lon float64, taxiType string) ([]NearbyDriver, error) {
	drivers, err := s.repo.FindByTaxiType(ctx, taxiType)
	if err != nil {
		return nil, err
	}

	const radiusKm = 6.0
	var nearby []NearbyDriver

	for _, d := range drivers {
		dist := haversineKm(lat, lon, d.Location.Lat, d.Location.Lon)
		if dist <= radiusKm {
			nearby = append(nearby, NearbyDriver{
				FirstName:  d.FirstName,
				LastName:   d.LastName,
				Plate:      d.Plate,
				DistanceKm: dist,
			})
		}
	}

	sort.Slice(nearby, func(i, j int) bool {
		return nearby[i].DistanceKm < nearby[j].DistanceKm
	})

	return nearby, nil
}

func haversineKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	toRad := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)

	lat1Rad := toRad(lat1)
	lat2Rad := toRad(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

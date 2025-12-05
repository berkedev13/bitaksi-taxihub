package driver

import (
	"context"
)

type Service interface {
	CreateDriver(ctx context.Context, req CreateDriverRequest) (*Driver, error)
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

package driver

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockRepository struct {
	InsertDriverFunc   func(ctx context.Context, d *Driver) (primitive.ObjectID, error)
	FindByIDFunc       func(ctx context.Context, id primitive.ObjectID) (*Driver, error)
	UpdateDriverFunc   func(ctx context.Context, d *Driver) error
	ListDriversFunc    func(ctx context.Context, page, pageSize int64) ([]Driver, error)
	FindByTaxiTypeFunc func(ctx context.Context, taxiType string) ([]Driver, error)
}

func (m *mockRepository) InsertDriver(ctx context.Context, d *Driver) (primitive.ObjectID, error) {
	if m.InsertDriverFunc != nil {
		return m.InsertDriverFunc(ctx, d)
	}
	return primitive.NilObjectID, nil
}

func (m *mockRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Driver, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) UpdateDriver(ctx context.Context, d *Driver) error {
	if m.UpdateDriverFunc != nil {
		return m.UpdateDriverFunc(ctx, d)
	}
	return nil
}

func (m *mockRepository) ListDrivers(ctx context.Context, page, pageSize int64) ([]Driver, error) {
	if m.ListDriversFunc != nil {
		return m.ListDriversFunc(ctx, page, pageSize)
	}
	return nil, nil
}

func (m *mockRepository) FindByTaxiType(ctx context.Context, taxiType string) ([]Driver, error) {
	if m.FindByTaxiTypeFunc != nil {
		return m.FindByTaxiTypeFunc(ctx, taxiType)
	}
	return nil, nil
}

func TestCreateDriver_Success(t *testing.T) {
	mockID := primitive.NewObjectID()
	mockRepo := &mockRepository{
		InsertDriverFunc: func(ctx context.Context, d *Driver) (primitive.ObjectID, error) {
			if d.FirstName == "" || d.Plate == "" {
				t.Errorf("expected driver fields to be set")
			}
			return mockID, nil
		},
	}

	svc := NewService(mockRepo)

	req := CreateDriverRequest{
		FirstName: "Ahmet",
		LastName:  "Demir",
		Plate:     "34ABC123",
		TaxiType:  "sari",
		CarBrand:  "Toyota",
		CarModel:  "Corolla",
		Lat:       41.0,
		Lon:       29.0,
	}

	ctx := context.Background()
	d, err := svc.CreateDriver(ctx, req)
	if err != nil {
		t.Fatalf("CreateDriver returned error: %v", err)
	}
	if d == nil {
		t.Fatalf("expected driver, got nil")
	}
	if d.ID != mockID {
		t.Errorf("expected ID %v, got %v", mockID, d.ID)
	}
	if d.FirstName != req.FirstName {
		t.Errorf("expected FirstName %q, got %q", req.FirstName, d.FirstName)
	}
}

func TestUpdateDriver_NotFound(t *testing.T) {
	mockRepo := &mockRepository{
		FindByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*Driver, error) {
			return nil, nil
		},
	}

	svc := NewService(mockRepo)

	ctx := context.Background()
	_, err := svc.UpdateDriver(ctx, primitive.NewObjectID().Hex(), UpdateDriverRequest{})
	if !errors.Is(err, ErrDriverNotFound) {
		t.Fatalf("expected ErrDriverNotFound, got %v", err)
	}
}

func TestUpdateDriver_PartialUpdate(t *testing.T) {
	existing := &Driver{
		ID:        primitive.NewObjectID(),
		FirstName: "Ahmet",
		LastName:  "Demir",
		Plate:     "34ABC123",
		TaxiType:  "sari",
		CarBrand:  "Toyota",
		CarModel:  "Corolla",
		Location: Location{
			Lat: 41.0,
			Lon: 29.0,
		},
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	var updatedSaved *Driver

	mockRepo := &mockRepository{
		FindByIDFunc: func(ctx context.Context, id primitive.ObjectID) (*Driver, error) {
			return existing, nil
		},
		UpdateDriverFunc: func(ctx context.Context, d *Driver) error {
			updatedSaved = d
			return nil
		},
	}

	svc := NewService(mockRepo)

	newBrand := "Hyundai"
	req := UpdateDriverRequest{
		CarBrand: &newBrand,
	}

	ctx := context.Background()
	updated, err := svc.UpdateDriver(ctx, existing.ID.Hex(), req)
	if err != nil {
		t.Fatalf("UpdateDriver returned error: %v", err)
	}

	if updated.CarBrand != newBrand {
		t.Errorf("expected CarBrand %q, got %q", newBrand, updated.CarBrand)
	}
	if updated.FirstName != existing.FirstName {
		t.Errorf("expected FirstName unchanged, got %q", updated.FirstName)
	}
	if updatedSaved == nil {
		t.Fatalf("expected UpdateDriver to be called on repo")
	}
	if updatedSaved.UpdatedAt.Before(existing.UpdatedAt) {
		t.Errorf("expected UpdatedAt to be refreshed")
	}
}

func TestGetNearbyDrivers_FilterAndSort(t *testing.T) {
	centerLat := 41.0
	centerLon := 29.0

	drivers := []Driver{
		{
			FirstName: "Yakın1",
			Plate:     "1",
			Location:  Location{Lat: 41.0001, Lon: 29.0001},
		},
		{
			FirstName: "Uzak",
			Plate:     "2",
			Location:  Location{Lat: 42.0, Lon: 30.0},
		},
		{
			FirstName: "Yakın2",
			Plate:     "3",
			Location:  Location{Lat: 41.01, Lon: 29.01},
		},
	}

	mockRepo := &mockRepository{
		FindByTaxiTypeFunc: func(ctx context.Context, taxiType string) ([]Driver, error) {
			if taxiType != "sari" {
				t.Errorf("expected taxiType 'sari', got %q", taxiType)
			}
			return drivers, nil
		},
	}

	svc := NewService(mockRepo)

	ctx := context.Background()
	result, err := svc.GetNearbyDrivers(ctx, centerLat, centerLon, "sari")
	if err != nil {
		t.Fatalf("GetNearbyDrivers returned error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 nearby drivers, got %d", len(result))
	}

	if result[0].FirstName != "Yakın1" {
		t.Errorf("expected Yakın1 to be first, got %s", result[0].FirstName)
	}
	if result[1].FirstName != "Yakın2" {
		t.Errorf("expected Yakın2 to be second, got %s", result[1].FirstName)
	}
}

func TestHaversineKm_ZeroDistance(t *testing.T) {
	lat := 41.0
	lon := 29.0

	dist := haversineKm(lat, lon, lat, lon)
	if dist > 0.0001 {
		t.Errorf("expected distance ~0, got %f", dist)
	}
}

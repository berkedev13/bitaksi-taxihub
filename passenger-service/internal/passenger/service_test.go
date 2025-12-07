package passenger

import (
	"context"
	"math"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mockRepo struct {
	passengers map[primitive.ObjectID]*Passenger
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		passengers: make(map[primitive.ObjectID]*Passenger),
	}
}

func (m *mockRepo) Insert(ctx context.Context, p *Passenger) (primitive.ObjectID, error) {
	if p.ID.IsZero() {
		p.ID = primitive.NewObjectID()
	}
	m.passengers[p.ID] = p
	return p.ID, nil
}

func (m *mockRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*Passenger, error) {
	p, ok := m.passengers[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}

func (m *mockRepo) Update(ctx context.Context, p *Passenger) error {
	m.passengers[p.ID] = p
	return nil
}

func (m *mockRepo) List(ctx context.Context, page, pageSize int64) ([]Passenger, error) {
	var all []Passenger
	for _, p := range m.passengers {
		all = append(all, *p)
	}

	start := (page - 1) * pageSize
	if start >= int64(len(all)) {
		return []Passenger{}, nil
	}

	end := start + pageSize
	if end > int64(len(all)) {
		end = int64(len(all))
	}

	return all[start:end], nil
}

func (m *mockRepo) FindAll(ctx context.Context) ([]Passenger, error) {
	var all []Passenger
	for _, p := range m.passengers {
		all = append(all, *p)
	}
	return all, nil
}

// --- Tests ---

func TestService_Create(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	req := CreatePassengerRequest{
		FirstName: "Mehmet",
		LastName:  "Yılmaz",
		Phone:     "+905551112233",
		Lat:       41.0431,
		Lon:       29.0099,
	}

	ctx := context.Background()
	p, err := svc.Create(ctx, req)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if p.ID.IsZero() {
		t.Fatalf("expected non-zero ID")
	}

	if p.FirstName != req.FirstName || p.LastName != req.LastName || p.Phone != req.Phone {
		t.Errorf("fields not set correctly, got %+v", p)
	}

	if p.Location.Lat != req.Lat || p.Location.Lon != req.Lon {
		t.Errorf("location not set correctly, got %+v", p.Location)
	}

	if p.CreatedAt.IsZero() || p.UpdatedAt.IsZero() {
		t.Errorf("expected timestamps to be set, got createdAt=%v updatedAt=%v", p.CreatedAt, p.UpdatedAt)
	}
}

func TestService_Update_Success(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	initial := &Passenger{
		ID:        primitive.NewObjectID(),
		FirstName: "Mehmet",
		LastName:  "Yılmaz",
		Phone:     "+905551112233",
		Location: Location{
			Lat: 41.0431,
			Lon: 29.0099,
		},
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}
	repo.passengers[initial.ID] = initial

	oldUpdatedAt := initial.UpdatedAt

	newPhone := "+905559998877"
	newLat := 41.05

	req := UpdatePassengerRequest{
		Phone: &newPhone,
		Lat:   &newLat,
	}

	ctx := context.Background()
	updated, err := svc.Update(ctx, initial.ID.Hex(), req)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if updated.Phone != newPhone {
		t.Errorf("expected phone %s, got %s", newPhone, updated.Phone)
	}
	if updated.Location.Lat != newLat {
		t.Errorf("expected lat %.4f, got %.4f", newLat, updated.Location.Lat)
	}

	if !updated.UpdatedAt.After(oldUpdatedAt) {
		t.Errorf("expected UpdatedAt to be refreshed, old=%v new=%v", oldUpdatedAt, updated.UpdatedAt)
	}
}

func TestService_Update_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	randomID := primitive.NewObjectID().Hex()
	req := UpdatePassengerRequest{
		Phone: strPtr("+905500000000"),
	}

	ctx := context.Background()
	_, err := svc.Update(ctx, randomID, req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != ErrPassengerNotFound {
		t.Fatalf("expected ErrPassengerNotFound, got %v", err)
	}
}

func TestService_List(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	p1 := &Passenger{
		ID:        primitive.NewObjectID(),
		FirstName: "A",
		LastName:  "One",
	}
	p2 := &Passenger{
		ID:        primitive.NewObjectID(),
		FirstName: "B",
		LastName:  "Two",
	}
	repo.passengers[p1.ID] = p1
	repo.passengers[p2.ID] = p2

	ctx := context.Background()
	list, err := svc.List(ctx, 1, 10)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 passengers, got %d", len(list))
	}
}

func TestService_GetNearby(t *testing.T) {
	repo := newMockRepo()
	svc := NewService(repo)

	centerLat, centerLon := 41.0431, 29.0099

	near := &Passenger{
		ID:        primitive.NewObjectID(),
		FirstName: "Near",
		Location: Location{
			Lat: centerLat,
			Lon: centerLon,
		},
	}

	far := &Passenger{
		ID:        primitive.NewObjectID(),
		FirstName: "Far",
		Location: Location{
			Lat: 40.0,
			Lon: 28.0,
		},
	}

	repo.passengers[near.ID] = near
	repo.passengers[far.ID] = far

	ctx := context.Background()
	result, err := svc.GetNearby(ctx, centerLat, centerLon)
	if err != nil {
		t.Fatalf("GetNearby returned error: %v", err)
	}

	foundNear := false
	foundFar := false

	for _, p := range result {
		if p.FirstName == "Near" {
			foundNear = true
		}
		if p.FirstName == "Far" {
			foundFar = true
		}
	}

	if !foundNear {
		t.Errorf("expected to find Near passenger in nearby result")
	}
	if foundFar {
		t.Errorf("did not expect to find Far passenger in nearby result")
	}
}

func TestHaversineKm_ZeroDistance(t *testing.T) {
	lat, lon := 41.0431, 29.0099
	d := haversineKm(lat, lon, lat, lon)
	if math.Abs(d) > 0.0001 {
		t.Errorf("expected distance ~0, got %f", d)
	}
}

func strPtr(s string) *string {
	return &s
}

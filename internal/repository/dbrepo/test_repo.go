package dbrepo

import (
	"context"
	"errors"
	"github.com/Deeksharma/bookings/internal/models"
	"time"
)

// implements all the functions of DatabaseRepo

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into DB
func (m *testDBRepo) InsertReservation(ctx context.Context, res models.Reservation) (int, error) {
	if res.LastName == "Loe" {
		return 0, errors.New("insert failed")
	}
	return 1, nil
}

// we would need to create models for our application to understand the db schema - dao - entity

// InsertRoomRestriction inserts a room restriction into DB
func (m *testDBRepo) InsertRoomRestriction(ctx context.Context, res models.RoomRestriction) error {
	if res.RoomId == 2 {
		return errors.New("insert failed")
	}
	return nil
}

// SearchAvailabilityForRoomByDates returns true if availability exists for roomId and returns false if roomId is not available
func (m *testDBRepo) SearchAvailabilityForRoomByDates(ctx context.Context, start, end time.Time, roomId int) (bool, error) {
	return true, nil
}

// SearchAvailabilityForAllRoomsByDates returns a slice of available rooms if any for given date range
func (m *testDBRepo) SearchAvailabilityForAllRoomsByDates(ctx context.Context, start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomById gives the room object by room id
func (m *testDBRepo) GetRoomById(ctx context.Context, id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("error")
	}
	room.ID = id
	return room, nil
}

func (m *testDBRepo) GetUserByID(ctx context.Context, id int) (models.User, error) {
	return models.User{}, nil
}

func (m *testDBRepo) UpdateUser(ctx context.Context, user models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(ctx context.Context, email, textPassword string) (int, string, error) {
	return 0, "", nil
}

func (m *testDBRepo) AllReservations(ctx context.Context) ([]models.Reservation, error) {
	return []models.Reservation{}, nil
}

func (m *testDBRepo) AllNewReservations(ctx context.Context) ([]models.Reservation, error) {
	return []models.Reservation{}, nil
}

func (m *testDBRepo) GetReservationById(ctx context.Context, id int) (models.Reservation, error) {
	return models.Reservation{}, nil
}

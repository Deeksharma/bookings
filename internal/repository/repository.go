package repository

import (
	"context"
	"github.com/Deeksharma/bookings/internal/models"
	"time"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(ctx context.Context, res models.Reservation) (int, error)
	InsertRoomRestriction(ctx context.Context, res models.RoomRestriction) error
	SearchAvailabilityForRoomByDates(ctx context.Context, start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRoomsByDates(ctx context.Context, start, end time.Time) ([]models.Room, error)
	GetRoomById(ctx context.Context, id int) (models.Room, error)
	GetUserByID(ctx context.Context, id int) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) error
	Authenticate(ctx context.Context, email, textPassword string) (int, string, error)
	AllReservations(ctx context.Context) ([]models.Reservation, error)
	AllNewReservations(ctx context.Context) ([]models.Reservation, error)
	GetReservationById(ctx context.Context, id int) (models.Reservation, error)
}

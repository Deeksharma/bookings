package dbrepo

import (
	"context"
	"errors"
	"github.com/Deeksharma/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// implements all the functions of DatabaseRepo

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into DB
func (m *postgresDBRepo) InsertReservation(ctx context.Context, res models.Reservation) (int, error) {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var newId int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
				end_date, room_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	err := m.DB.QueryRow(nctx, stmt, res.FirstName, res.LastName, res.Email,
		res.Phone, res.StartDate, res.EndDate,
		res.RoomId, time.Now(), time.Now()).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

// we would need to create models for our application to understand the db schema - dao - entity

// InsertRoomRestriction inserts a room restriction into DB
func (m *postgresDBRepo) InsertRoomRestriction(ctx context.Context, res models.RoomRestriction) error {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date,
				end_date, room_id, reservation_id, restriction_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := m.DB.Exec(nctx,
		stmt,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		res.ReservationId,
		res.RestrictionId,
		time.Now(),
		time.Now())

	if err != nil {
		return err
	}
	return nil
}

// SearchAvailabilityForRoomByDates returns true if availability exists for roomId and returns false if roomId is not available
func (m *postgresDBRepo) SearchAvailabilityForRoomByDates(ctx context.Context, start, end time.Time, roomId int) (bool, error) {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	stmt := `select count(*)
		from room_restrictions x
		where room_id = $1 and start_date < $2 and end_date > $3`

	var numRows int
	row := m.DB.QueryRow(nctx, stmt, roomId, end, start)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	return numRows == 0, nil
}

// SearchAvailabilityForAllRoomsByDates returns a slice of available rooms if any for given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRoomsByDates(ctx context.Context, start, end time.Time) ([]models.Room, error) {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	stmt := `select r.id, r.room_name from rooms as r
		where r.id not in 
		(select rr.room_id from room_restrictions rr where rr.start_date <= $1 and rr.end_date >= $2)`

	var rooms []models.Room
	rows, err := m.DB.Query(nctx, stmt, end, start)

	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}
	defer rows.Close()
	return rooms, nil
}

// GetRoomById gives the room object by room id
func (m *postgresDBRepo) GetRoomById(ctx context.Context, id int) (models.Room, error) {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var room models.Room

	stmt := `select id, room_name, created_at, updated_at from rooms where id = $1`

	row := m.DB.QueryRow(nctx, stmt, id)

	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt)

	if err != nil {
		return room, err
	}
	return room, nil
}

// GetUserByID gets user by id
func (m *postgresDBRepo) GetUserByID(ctx context.Context, id int) (models.User, error) {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var user models.User

	stmt := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where id = $1`

	row := m.DB.QueryRow(nctx, stmt, id)
	err := row.Scan(&user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.AccessLevel,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

// UpdateUser updates user in the database
func (m *postgresDBRepo) UpdateUser(ctx context.Context, user models.User) error {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `update users set first_name= $1, last_name= $2, email=$3, access_level= $4, updated_at= $5` // for password we'll use different process

	_, err := m.DB.Exec(nctx, stmt, user.FirstName, user.LastName, user.Email, user.AccessLevel, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// we are going to store password using hash

// Authenticate authenticates a user and returns id and hashed password of user / error
func (m *postgresDBRepo) Authenticate(ctx context.Context, email, textPassword string) (int, string, error) {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	stmt := "select id, password from users where email = $1"
	row := m.DB.QueryRow(nctx, stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(textPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// AllReservations returns a slice of all reservations
func (m *postgresDBRepo) AllReservations(ctx context.Context) ([]models.Reservation, error) {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	stmt := "select r.id, r.first_name, r.last_name, r.phone, r.start_date, r.end_date," +
		" r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name from reservations r left join" +
		" rooms rm on (r.room_id = rm.id) order by r.start_date asc"

	rows, err := m.DB.Query(nctx, stmt)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomId,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Processed,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, reservation)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	defer rows.Close()
	return reservations, err
}

// AllNewReservations returns a slice of all reservations
func (m *postgresDBRepo) AllNewReservations(ctx context.Context) ([]models.Reservation, error) {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	stmt := "select r.id, r.first_name, r.last_name, r.phone, r.start_date, r.end_date," +
		" r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name from reservations r left join" +
		" rooms rm on (r.room_id = rm.id) where r.processed = 0 order by r.start_date asc"

	rows, err := m.DB.Query(nctx, stmt)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(&reservation.ID,
			&reservation.FirstName,
			&reservation.LastName,
			&reservation.Phone,
			&reservation.StartDate,
			&reservation.EndDate,
			&reservation.RoomId,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
			&reservation.Processed,
			&reservation.Room.ID,
			&reservation.Room.RoomName,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, reservation)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}
	defer rows.Close()
	return reservations, err
}

// GetReservationById returns a reservation for the respective id
func (m *postgresDBRepo) GetReservationById(ctx context.Context, id int) (models.Reservation, error) {
	nctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var reservation models.Reservation

	stmt := "select r.id, r.first_name, r.last_name, r.phone, r.email, r.start_date, r.end_date," +
		" r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name from reservations r left join" +
		" rooms rm on (r.room_id = rm.id) where r.id = $1"

	row := m.DB.QueryRow(nctx, stmt, id)

	err := row.Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Phone,
		&reservation.Email,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomId,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Processed,
		&reservation.Room.ID,
		&reservation.Room.RoomName,
	)

	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

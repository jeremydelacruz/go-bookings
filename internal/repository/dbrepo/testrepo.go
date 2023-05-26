package dbrepo

import (
	"errors"
	"time"

	"github.com/jeremydelacruz/go-bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// induce error for testing
	if res.RoomID == 999 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	// induce error for testing
	if r.RoomID == 1000 {
		return errors.New("some error")
	}
	return nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	return true, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room

	// induce error for testing
	if id == 999 {
		return room, errors.New("some error")
	}

	return room, nil
}

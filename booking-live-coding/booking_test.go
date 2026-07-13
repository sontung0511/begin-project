package booking_live_coding

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCancelBooking(t *testing.T) {
	db := setupTestDB(t)
	service := NewService(db)
	event := Event{
		ID:               uuid.NewString(),
		AvailableTickets: 100,
	}
	require.NoError(t, db.Create(&event).Error)
	booking := Booking{
		ID:       uuid.NewString(),
		EventID:  event.ID,
		Quantity: 1,
		Status:   "PENDING",
	}
	require.NoError(t, db.Create(&booking).Error)
	const workerCount = 2
	start := make(chan struct{})
	results := make(chan error, workerCount)
	var wg sync.WaitGroup
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			<-start
			results <- service.CancelBooking(context.Background(), booking.ID)
		}()
	}
	close(start)
	wg.Wait()
	close(results)

	successCount := 0
	notPendingCount := 0
	for e := range results {
		switch {
		case e == nil:
			successCount++
		case errors.Is(e, ErrBookingNotPending):
			notPendingCount++
		default:
			t.Fatal("unexpected error", e)

		}
	}
	require.Equal(t, workerCount, successCount)
	require.Equal(t, workerCount, notPendingCount)
	var actualBooking Booking
	require.NoError(t, db.First(&actualBooking, "id = ?", booking.ID).Error)
	require.Equal(t, "CANCELLED", actualBooking.Status)

	var actualEvent Event
	require.NoError(t,
		db.First(&actualEvent, "id = ?", event.ID).Error,
	)

	require.Equal(t, 100, actualEvent.AvailableTickets)
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=ticket_booking_test port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	sqlDB, err := db.DB()
	require.NoError(t, err)

	require.NoError(t, sqlDB.Ping())

	require.NoError(t, db.AutoMigrate(
		&Event{},
		&Booking{},
	))

	require.NoError(t, db.Exec(`
		TRUNCATE TABLE bookings, events
		RESTART IDENTITY CASCADE
	`).Error)

	t.Cleanup(func() {
		_ = db.Exec(`
			TRUNCATE TABLE bookings, events
			RESTART IDENTITY CASCADE
		`).Error

		_ = sqlDB.Close()
	})

	return db
}

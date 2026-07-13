package booking_live_coding

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrBookingNotFound   = errors.New("booking not found")
	ErrBookingNotPending = errors.New("booking is not pending")
)

type Booking struct {
	ID       string
	EventID  string
	Quantity int
	Status   string
}

type Event struct {
	ID               string
	AvailableTickets int
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CancelBooking requirements:
//
// 1. Chỉ được cancel booking có status PENDING.
// 2. Chuyển status từ PENDING sang CANCELLED.
// 3. Hoàn Quantity vào events.available_tickets.
// 4. Booking update và ticket release phải cùng transaction.
// 5. Hai request chạy đồng thời chỉ được hoàn vé đúng một lần.
// 6. Trả ErrBookingNotFound nếu booking không tồn tại.
// 7. Trả ErrBookingNotPending nếu booking đã CONFIRMED hoặc CANCELLED.
func (s *Service) CancelBooking(
	ctx context.Context,
	bookingID string,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var booking Booking

		err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&booking, "id = ?", bookingID).
			Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBookingNotFound
			}
			return err
		}

		if booking.Status != "PENDING" {
			return ErrBookingNotPending
		}

		result := tx.Model(&Booking{}).
			Where("id = ? AND status = ?", bookingID, "PENDING").
			Update("status", "CANCELLED")

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return ErrBookingNotPending
		}

		result = tx.Model(&Event{}).
			Where("id = ?", booking.EventID).
			UpdateColumn(
				"available_tickets",
				gorm.Expr("available_tickets + ?", booking.Quantity),
			)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("event not found")
		}

		return nil
	})
}

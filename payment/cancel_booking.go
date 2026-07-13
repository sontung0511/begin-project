package main

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Service) CancelBooking(ctx context.Context, bookingId int) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var booking Booking
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", bookingId).First(&booking).Error; err != nil {
			return err
		}
		if booking.Status != StatusPending {
			return nil
		}
		result := tx.Model(&Booking{}).Where("id = ? AND status = ?", bookingId, StatusPending).
			Update("status", StatusCancelled)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("Booking with id %d not found", bookingId)
		}
		resultEvent := tx.Model(&Event{}).Where("id = ?", booking.EventID).UpdateColumn("available_tickets", gorm.Expr("available_tickets + ?", booking.Quantity))
		if resultEvent.Error != nil {
			return resultEvent.Error
		}
		if resultEvent.RowsAffected == 0 {
			return fmt.Errorf("Not update event to quantity")
		}
		return nil
	})
}

package main

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Payment struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	BookingID int    `json:"booking_id"`
}
type Booking struct {
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Quantity int    `json:"quantity"`
	EventID  int    `json:"event_id"`
}
type Event struct {
	ID               int `json:"id"`
	AvailableTickets int `json:"available_tickets"`
}

const (
	StatusCancelled = "cancelled"
	StatusSuccess   = "success"
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) PaymentCallBack(ctx context.Context, paymentId int, callBackStatus string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//validate payment
		var payment Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&payment, "id = ?", paymentId).Error; err != nil {
			return err
		}
		if payment.Status != StatusPending {
			return nil
		}
		result := tx.Model(&Payment{}).
			Where("id = ? AND status = ?", paymentId, StatusPending).
			Update("status", callBackStatus)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("no such payment")
		}
		if callBackStatus == StatusCancelled {
			var booking Booking
			resultBooking := tx.Model(booking).Where("id = ? AND status = ?", payment.BookingID, StatusPending).Update("status", StatusCancelled)
			if resultBooking.Error != nil {
				return resultBooking.Error
			}
			if resultBooking.RowsAffected == 0 {
				return errors.New("booking not update error")
			}
			resultEvent := tx.Model(&Event{}).Where("id = ? AND available_tickets > ?", booking.EventID, booking.Quantity).
				Update("available_tickets + ?", booking.Quantity)
			if resultEvent.Error != nil {
				return resultBooking.Error
			}
			if resultEvent.RowsAffected == 0 {
				return errors.New("event not update error")
			}
		}
		if callBackStatus == StatusSuccess {
			resultBooking := tx.Model(&Booking{}).
				Where("id = ? AND status = ?", payment.BookingID, StatusPending).
				Update("status", StatusConfirmed)
			if resultBooking.Error != nil {
				return resultBooking.Error
			}
			if resultBooking.RowsAffected == 0 {
				return errors.New("booking not update error")
			}
		}
		return nil
	})
}

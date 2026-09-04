package pos

import (
	"database/sql"
	"errors"
	"fmt"
)

// ErrEventStarted marks a data-reset attempt after Start Event is set
// (ADR-0007). It maps to a locked response instead of a wipe.
var ErrEventStarted = errors.New("event started")

// EventStarted reports whether the System Admin has set Start Event, the
// one-way flag that locks the data-reset tool out for the rest of the event.
func (service *OrderService) EventStarted() (bool, error) {
	var value string
	err := service.DB.QueryRow(`SELECT value FROM metadata WHERE key = 'event_started'`).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read event_started: %w", err)
	}
	return value == "true", nil
}

// StartEvent sets the Start Event flag (ADR-0007). It cannot be unset.
func (service *OrderService) StartEvent() error {
	_, err := service.DB.Exec(
		`INSERT INTO metadata (key, value) VALUES ('event_started', 'true')
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`)
	if err != nil {
		return fmt.Errorf("set event_started: %w", err)
	}
	return nil
}

// ResetAllOrders hard-deletes every order and restarts the order-number
// counter at StartingOrderNumber (ADR-0007). It refuses once Start Event is
// set, so a real event day's sales can never be erased by mistake.
func (service *OrderService) ResetAllOrders() error {
	started, err := service.EventStarted()
	if err != nil {
		return err
	}
	if started {
		return ErrEventStarted
	}

	transaction, err := service.DB.Begin()
	if err != nil {
		return fmt.Errorf("begin reset transaction: %w", err)
	}
	defer transaction.Rollback()

	if _, err := transaction.Exec(`DELETE FROM transactions`); err != nil {
		return fmt.Errorf("delete transactions: %w", err)
	}
	if _, err := transaction.Exec(`DELETE FROM metadata WHERE key IN ('business_date', 'next_order_number')`); err != nil {
		return fmt.Errorf("delete order-number metadata: %w", err)
	}
	return transaction.Commit()
}

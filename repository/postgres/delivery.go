package repository

import (
	"database/sql"
	"time"

	"github.com/allisson/hammer"
	"github.com/jmoiron/sqlx"
)

// Delivery is a implementation of hammer.DeliveryRepository
type Delivery struct {
	db *sqlx.DB
}

// Find returns hammer.Delivery by id
func (d *Delivery) Find(id string) (hammer.Delivery, error) {
	delivery := hammer.Delivery{}
	err := d.db.Get(&delivery, sqlDeliveryFind, id)
	return delivery, err
}

// FindAll returns []hammer.Delivery by limit and offset
func (d *Delivery) FindAll(limit, offset int) ([]hammer.Delivery, error) {
	deliveries := []hammer.Delivery{}
	err := d.db.Select(&deliveries, sqlDeliveryFindAll, limit, offset)
	return deliveries, err
}

// FindToDispatch returns []hammer.Delivery ready to dispatch by limit and offset
func (d *Delivery) FindToDispatch(limit, offset int) ([]string, error) {
	deliveries := []string{}
	status := "pending"
	now := time.Now().UTC()
	err := d.db.Select(&deliveries, sqlDeliveryFindToDispatch, status, now, limit, offset)
	return deliveries, err
}

// Store a hammer.Delivery on database (create or update)
func (d *Delivery) Store(tx hammer.TxRepository, delivery *hammer.Delivery) error {
	_, err := d.Find(delivery.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return tx.Exec(sqlDeliveryCreate, delivery)
		}
		return err
	}
	return tx.Exec(sqlDeliveryUpdate, delivery)
}

// NewDelivery returns a new Delivery with db connection
func NewDelivery(db *sqlx.DB) Delivery {
	return Delivery{db: db}
}

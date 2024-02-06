package model

import "time"

// Item represents a simple bit of data stored in the database.
type Item struct {
	// Internal ID in the database.
	ID int

	// Title of the item.
	Title string

	// Description of the item.
	Description string

	// Time the item was created.
	CreationDate time.Time
}

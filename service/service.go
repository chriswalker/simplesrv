// Package service contains the primary business logic for the server.
package service

import (
	"log/slog"

	"github.com/chriswalker/simplesrv/db"
	"github.com/chriswalker/simplesrv/model"
)

type itemGetter interface {
	GetItems() ([]model.Item, error)
}

// ItemService encapsulates the business logic of the item server.
type ItemService struct {
	db     itemGetter
	logger *slog.Logger
}

// NewItemService creates a new service object, complete
// with configured SQLite database.
func NewItemService(filename string, logger *slog.Logger) (*ItemService, error) {
	s := &ItemService{
		logger: logger,
	}
	d, err := db.New(filename)
	if err != nil {
		return nil, err
	}
	s.db = d

	return s, nil
}

// GetItems returns all items from the database.
func (s *ItemService) GetItems() ([]model.Item, error) {
	return s.db.GetItems()
}

package model

import "context"

type [Entity] struct {
	ID      string `db:"_id"`
	UUID    string `db:"uuid"`
	AuditableEntity
}

type SQLStore interface {
	Create(ctx context.Context, model *) (*, error)
	Update(ctx context.Context, model *) (*, error)
	Delete(ctx context.Context, model *) error
	FindByID(ctx context.Context, id string) (*, error)
	FindAll(ctx context.Context) ([]*, error)
}

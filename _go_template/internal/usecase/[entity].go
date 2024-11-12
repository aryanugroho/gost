package usecase

import (
	"context"
	"errors"

	"github.com/aryanugroho//internal/model"
)

type Payload struct {
	UUID    string `json:"uuid"`
}

type [Entity] struct {
	ID      string `json:"id"`
	UUID    string `json:"uuid"`
	model.AuditableEntity
}


func (u *[Entity]) Create(ctx context.Context, payload *Payload) (*, error) {
	ctx = u.Infrastructure.SQLStore().BeginTx(ctx)
	defer u.Infrastructure.SQLStore().RollbackTx(ctx)

	, err := u.Infrastructure.SQLStore().Store().Create(ctx, &model.{
		UUID:    payload.UUID,
	})
	if err != nil {
		return nil, errors.New("failed to create post")
	}

	err = u.Infrastructure.SQLStore().CommitTx(ctx)
	if err != nil {
		return nil, errors.New("failed to create post")
	}

	return &{
		ID:              .ID,
		UUID:            .UUID,
		AuditableEntity: .AuditableEntity,
	}, nil
}

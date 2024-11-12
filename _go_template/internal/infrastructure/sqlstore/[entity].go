package sqlstore

import (
	"context"

	"github.com/aryanugroho//pkg/db"
	"github.com/aryanugroho//internal/model"
)

// implementation of model store
type Store struct {
	master *db.SqlxDBWrapper
	slave  *db.SqlxDBWrapper
}

func NewStore(master *db.SqlxDBWrapper, slave *db.SqlxDBWrapper) *Store {
	return &Store{
		master: master,
		slave:  slave,
	}
}

func (u *Store) Create(ctx context.Context, model *model.) (*model., error) {
	const createSQL  = `INSERT INTO s (uuid, created_at, updated_at) VALUES (?, ?, ?)`

	dbs := u.master
	tx, ok := ctx.Value(Tx).(*db.SqlxDBWrapper)
	if ok {
		dbs = tx
	}

	_, err := dbs.ExecContext(ctx, createSQL, model.UUID, model.CreatedAt, model.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (u *Store) Update(ctx context.Context, model *model.) (*model., error) {
	const updateSQL  = `UPDATE s SET uuid = ?, updated_at = ? WHERE id = ?`

	dbs := u.master
	tx, ok := ctx.Value(Tx).(*db.SqlxDBWrapper)
	if ok {
		dbs = tx
	}

	_, err := dbs.ExecContext(ctx, updateSQL, model.UUID,model.UpdatedAt, model.ID)
	if err != nil {
		return nil, err
	}

	return model, nil
}

func (u *Store) Delete(ctx context.Context, model *model.) error {
	const deleteSQL  = `DELETE FROM s WHERE id = ?`

	dbs := u.master
	tx, ok := ctx.Value(Tx).(*db.SqlxDBWrapper)
	if ok {
		dbs = tx
	}

	_, err := dbs.ExecContext(ctx, deleteSQL, model.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *Store) FindByID(ctx context.Context, id string) (*model., error) {
	findByIDSQL    = `SELECT uuid, created_at, updated_at FROM s WHERE id = ?`

	dbs := u.slave
	tx, ok := ctx.Value(Tx).(*db.SqlxDBWrapper)
	if ok {
		dbs = tx
	}

	var model []model.
	err := dbs.SelectContext(ctx, &model, findByIDSQL, id)
	if err != nil {
		return nil, err
	}

	if len(model) < 1 {
		return nil, nil
	}

	return &model[0], nil
}

func (u *Store) FindAll(ctx context.Context) ([]*model., error) {
	const findAllSQL = `SELECT uuid, created_at, updated_at FROM s ORDER BY created_at DESC`

	dbs := u.slave
	tx, ok := ctx.Value(Tx).(*db.SqlxDBWrapper)
	if ok {
		dbs = tx
	}

	var model []*model.
	err := dbs.SelectContext(ctx, model, findAllSQL)
	if err != nil {
		return nil, err
	}

	return model, nil
}

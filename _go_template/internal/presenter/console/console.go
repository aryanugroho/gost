package console

import "github.com/aryanugroho//internal/infrastructure/sqlstore"

type Console struct {
	Migrator *sqlstore.Migrator
}

func Init() (*Console, error) {
	sqlMigrator, err := sqlstore.NewMigrator()
	if err != nil {
		return nil, err
	}

	return &Console{
		Migrator: sqlMigrator,
	}, nil
}

package app

import (
	errors "github.com/Red-Sock/trace-errors"
	"github.com/godverv/matreshka-be/internal/clients/sqldb"
)

func (a *App) InitDataSources() (err error) {
	a.Sqlite, err = sqldb.New(a.Cfg.DataSources.Sqlite)
	if err != nil {
		return errors.Wrap(err, "error during sql connection initialization")
	}

	return nil
}

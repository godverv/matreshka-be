// Code generated by RedSock CLI. DO NOT EDIT.

package app

import (
	errors "go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka-be/internal/clients/sqldb"
)

func (a *App) InitDataSources() (err error) {
	a.Sqlite, err = sqldb.New(a.Cfg.DataSources.Sqlite)
	if err != nil {
		return errors.Wrap(err, "error during sql connection initialization")
	}

	return nil
}

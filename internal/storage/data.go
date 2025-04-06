package storage

import (
	"context"
	"database/sql"
	"time"

	"go.redsock.ru/evon"

	"go.vervstack.ru/matreshka/internal/domain"
)

type Data interface {
	// GetConfigNodes returns root node of parsed config
	GetConfigNodes(ctx context.Context, serviceName string, ver string) (*evon.Node, error)
	ListConfigs(ctx context.Context, req domain.ListConfigsRequest) (domain.ListConfigsResponse, error)

	Create(ctx context.Context, serviceConfig string) (int64, error)

	UpsertValues(ctx context.Context, req domain.PatchConfigRequest) error
	DeleteValues(ctx context.Context, req domain.PatchConfigRequest) error

	SetUpdatedAt(ctx context.Context, serviceName string, req time.Time) error

	Rename(ctx context.Context, oldName string, newName string) error

	WithTx(tx *sql.Tx) Data
}

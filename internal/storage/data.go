package storage

import (
	"context"
	"database/sql"

	"github.com/Red-Sock/evon"

	"github.com/godverv/matreshka-be/internal/domain"
)

type Data interface {
	// GetConfigNodes returns root node of parsed config
	GetConfigNodes(ctx context.Context, serviceName string) (*evon.Node, error)
	ListConfigs(ctx context.Context, req domain.ListConfigsRequest) ([]domain.ConfigListItem, error)

	Create(ctx context.Context, serviceConfig string) (int64, error)

	UpsertValues(ctx context.Context, serviceName string, req []domain.PatchConfig) error
	DeleteValues(ctx context.Context, serviceName string, req []domain.PatchConfig) error

	WithTx(tx *sql.Tx) Data
}
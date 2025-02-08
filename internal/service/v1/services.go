package v1

import (
	"go.vervstack.ru/matreshka-be/internal/service"
	"go.vervstack.ru/matreshka-be/internal/service/v1/config"
	"go.vervstack.ru/matreshka-be/internal/service/v1/subscription"
	"go.vervstack.ru/matreshka-be/internal/storage"
	"go.vervstack.ru/matreshka-be/internal/storage/tx_manager"
)

type Services struct {
	configService *config.CfgService
	pubSubService *subscription.PubSubService
}

func New(data storage.Data, txManager *tx_manager.TxManager) *Services {
	pubSubService := subscription.New()

	return &Services{
		configService: config.New(data, txManager, pubSubService),
		pubSubService: pubSubService,
	}
}

func (s *Services) ConfigService() service.ConfigService {
	return s.configService
}
func (s *Services) PubSubService() service.PubSubService {
	return s.pubSubService
}

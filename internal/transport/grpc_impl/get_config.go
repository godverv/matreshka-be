package grpc_impl

import (
	"context"

	"go.redsock.ru/evon"
	errors "go.redsock.ru/rerrors"
	"go.vervstack.ru/matreshka"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "go.vervstack.ru/matreshka-be/pkg/matreshka_be_api"
)

func (a *Impl) GetConfig(ctx context.Context, req *api.GetConfig_Request) (*api.GetConfig_Response, error) {
	cfgNodes, err := a.configService.GetNodes(ctx, req.GetServiceName())
	if err != nil {
		return nil, errors.Wrap(err)
	}

	targetConfig := matreshka.NewEmptyConfig()

	nodeStorage := evon.NodesToStorage([]*evon.Node{cfgNodes})

	err = evon.UnmarshalWithNodes(nodeStorage, &targetConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling config")
	}

	resp := &api.GetConfig_Response{}

	resp.Config, err = targetConfig.Marshal()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return resp, nil
}

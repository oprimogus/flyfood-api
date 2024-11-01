package adapter

import (
	"context"

	"github.com/oprimogus/cardapiogo/internal/services/adapter/identity"
	"github.com/oprimogus/cardapiogo/internal/services/adapter/storage"
	"github.com/oprimogus/cardapiogo/internal/services/aws"
	"github.com/oprimogus/cardapiogo/internal/services/keycloak"
)

type Factory interface {
	NewStorageService() storage.Service
	NewIdentityService() identity.Service
}

type ServiceFactory struct {
	storage          storage.Service
	identityProvider identity.Service
}

func (f ServiceFactory) NewStorageService() storage.Service {
	return f.storage
}

func (f ServiceFactory) NewIdentityService() identity.Service {
	return f.identityProvider
}

func NewServiceFactory() Factory {
	ctx := context.Background()
	awsInstance, err := aws.GetInstance(ctx)
	if err != nil {
		panic(err)
	}
	s3 := &awsInstance.S3

	k, err := keycloak.GetInstance()
	if err != nil {
		panic(err)
	}

	return ServiceFactory{
		storage:          s3,
		identityProvider: k,
	}
}

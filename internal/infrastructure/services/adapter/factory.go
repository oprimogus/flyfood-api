package adapter

import (
	"context"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/adapter/storage"
	"github.com/oprimogus/cardapiogo/internal/infrastructure/services/aws"
)

type Factory interface {
	NewStorageService() storage.Service
	//NewIdentityService() *keycloak.Service
}

type ServiceFactory struct {
	storage storage.Service
	//identity *keycloak.Service
}

func (f ServiceFactory) NewStorageService() storage.Service {
	return f.storage
}

//func (f ServiceFactory) NewIdentityService() *keycloak.Service {
//	return f.identity
//}

func NewServiceFactory() Factory {
	ctx := context.Background()
	awsInstance, err := aws.GetInstance(ctx)
	if err != nil {
		panic(err)
	}
	s3 := &awsInstance.S3

	//keycloakInstance, err := keycloak.GetInstance()
	//if err != nil {
	//	panic(err)
	//}

	return &ServiceFactory{
		storage: s3,
		//identity: nil,
	}
	//return &ServiceFactory{
	//	storage:  s3,
	//	identity: keycloakInstance,
	//}
}

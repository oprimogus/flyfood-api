package store

import (
	"context"
	"fmt"
	"mime/multipart"

	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

type useCaseSetProfileImage struct {
	repository Repository
}

func newUseCaseSetProfileImage(repository Repository) useCaseSetProfileImage {
	return useCaseSetProfileImage{
		repository: repository,
	}
}

func (c useCaseSetProfileImage) Execute(ctx context.Context, storeID string, image *multipart.FileHeader) (objectURL string, err error) {
	userID, ok := ctx.Value(string(logger.UserIDKey)).(string)
	if !ok {
		return "", fmt.Errorf("invalid userID: '%s'", userID)
	}
	if userID == "" {
		return "", fmt.Errorf("invalid userID: '%s'", userID)
	}

	isOwner, errIsOwner := c.repository.IsOwner(ctx, storeID, userID)
	if errIsOwner != nil {
		return "", errIsOwner
	}

	if !isOwner {
		return "", ErrNotOwner
	}

	url, errSaveProfileImage := c.repository.SetProfileImage(ctx, storeID, image)
	if errSaveProfileImage != nil {
		return "", errSaveProfileImage
	}

	return url, nil
}

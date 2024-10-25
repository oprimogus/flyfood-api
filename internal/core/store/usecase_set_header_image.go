package store

import (
	"context"
	"fmt"
	"mime/multipart"
)

type useCaseSetHeaderImage struct {
	repository Repository
}

func newUseCaseSetHeaderImage(repository Repository) useCaseSetHeaderImage {
	return useCaseSetHeaderImage{
		repository: repository,
	}
}

func (c useCaseSetHeaderImage) Execute(ctx context.Context, storeID string, image *multipart.FileHeader) (objectURL string, err error) {
	userID := ctx.Value("userID").(string)
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

	url, errSaveProfileImage := c.repository.SetHeaderImage(ctx, storeID, image)
	if errSaveProfileImage != nil {
		return "", errSaveProfileImage
	}

	return url, nil
}

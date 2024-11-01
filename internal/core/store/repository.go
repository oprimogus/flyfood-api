package store

import (
	"context"
	"mime/multipart"
)

type Repository interface {
	Create(ctx context.Context, params Store) (id string, err error)
	Update(ctx context.Context, userID string, params Store) error
	FindByID(ctx context.Context, id string) (Store, error)
	FindByFilter(ctx context.Context, params GetStoresFilterInput) (*[]Store, error)
	Delete(ctx context.Context, id string) error
	Enable(ctx context.Context, id string) error
	IsOwner(ctx context.Context, id, userID string) (bool, error)
	AddBusinessHour(ctx context.Context, storeID string, params []BusinessHours) error
	DeleteBusinessHour(ctx context.Context, storeID string, params []BusinessHours) error
	SetProfileImage(ctx context.Context, storeID string, image *multipart.FileHeader) (objectURL string, err error)
	SetHeaderImage(ctx context.Context, storeID string, image *multipart.FileHeader) (objectURL string, err error)
}

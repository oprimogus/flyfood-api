package store

import "github.com/oprimogus/flyfood-api/internal/infra/services/adapter/storage"

func GetBuckets() []storage.Bucket {
	return []storage.Bucket{
		ProfileBucket,
		HeaderBucket,
		ProductBucket,
	}
}

const (
	ProfileBucket storage.Bucket = "flyfood-store-profile-image"
	HeaderBucket  storage.Bucket = "flyfood-store-header-image"
	ProductBucket storage.Bucket = "flyfood-store-product-image"
)

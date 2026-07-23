package main

import (
	"context"
	"fmt"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/oprimogus/flyfood-api/internal/core/store"
	"github.com/oprimogus/flyfood-api/internal/infra/services/adapter"
	"github.com/oprimogus/flyfood-api/internal/infra/services/adapter/storage"
	"log"
)

func main() {
	ctx := context.Background()
	_ = config.Get()
	factory := adapter.NewServiceFactory()
	service := factory.NewStorageService()

	storeBuckets := store.GetBuckets()

	buckets := makeSlice(storeBuckets)

	for _, v := range buckets {
		existBucket, err := service.BucketExists(ctx, string(v))
		if err != nil {
			log.Fatalf("Não foi possível verificar se o bucket %s já existe ou não", string(v))
		}
		if existBucket {
			fmt.Printf("O bucket %s já existe\n", string(v))
		} else {
			err = service.CreateBucket(ctx, string(v))
			if err != nil {
				log.Fatalf("Não foi possível criar o bucket %s", string(v))
			}
			fmt.Printf("Bucket %s criado\n", string(v))
		}
	}
}

func calculateSum(b ...[]storage.Bucket) int {
	sum := 0
	bMap := make([][]storage.Bucket, len(b))

	_ = copy(bMap, b)

	for _, v := range bMap {
		sum += len(v)
	}
	return sum
}

func makeSlice(b ...[]storage.Bucket) []storage.Bucket {
	buckets := make([]storage.Bucket, calculateSum(b...))
	for _, v := range b {
		index := 0
		for _, k := range v {
			buckets[index] = k
			index++
		}
	}
	return buckets
}

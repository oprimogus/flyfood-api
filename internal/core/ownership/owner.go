package ownership

import (
	"time"

	"github.com/google/uuid"
)

type Owner struct {
	ID uuid.UUID
	SignatureActive bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
package model

import "github.com/google/uuid"

type PartInfo struct {
	PartUUID uuid.UUID
	Price    float64
}

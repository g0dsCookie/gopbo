package pbo

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidProductEntry = errors.New("invalid product entry")
	ErrFileCorrupted       = errors.New("file seems corrupted")
)

type InvalidPackingMethod struct {
	Packing PackingMethod
}

func (i *InvalidPackingMethod) Error() string {
	return fmt.Sprintf("invalid packing method %v", i.Packing)
}

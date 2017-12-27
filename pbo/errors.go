package pbo

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidProductEntry is returned when no valid product entry
	// could be found.
	ErrInvalidProductEntry = errors.New("invalid product entry")

	// ErrFileCorrupted is returned when the sha1 stored in the pbo doesn't matched
	// the calculated ones of this package.
	ErrFileCorrupted = errors.New("file seems corrupted")
)

// InvalidPackingMethod is returned when the packing method is unknown.
type InvalidPackingMethod struct {
	Packing PackingMethod // Packing contains the unknown packing method.
}

// Error returns a user readable error string.
func (i *InvalidPackingMethod) Error() string {
	return fmt.Sprintf("invalid packing method %x", i.Packing)
}

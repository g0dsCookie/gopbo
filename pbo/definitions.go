package pbo

// PackingMethod describes the packing method
// used to store the file in the pbo.
type PackingMethod uint32

const (
	// PackingMethodUncompressed defines the file as uncompressed
	PackingMethodUncompressed PackingMethod = 0x00000000

	// PackingMethodPacked defines the file as "packed" (compressed)
	PackingMethodPacked = 0x43707273

	// PackingMethodProductEntry defines the entry as product entry
	// and is usually found at the very start of the pbo file.
	PackingMethodProductEntry = 0x56657273
)

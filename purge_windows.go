package ryugraph

// PurgeHeap is a no-op on Windows — each DLL has its own heap and
// memory is reclaimed automatically when the DB handle is destroyed.
func PurgeHeap() {}

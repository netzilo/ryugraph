//go:build !cgo && (darwin || linux)

package ryugraph

// PurgeHeap is a no-op when CGo is disabled — malloc_trim/malloc_zone_pressure_relief
// are unavailable without CGo.
func PurgeHeap() {}

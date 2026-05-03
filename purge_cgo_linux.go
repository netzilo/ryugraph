//go:build cgo && linux

package ryugraph

/*
#include <malloc.h>
*/
import "C"

// PurgeHeap asks the C allocator to return free pages to the OS.
// Call after closing a large DB to recover RSS.
func PurgeHeap() {
	C.malloc_trim(0)
}

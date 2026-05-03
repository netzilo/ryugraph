//go:build !cgo

package ryugraph

import "errors"

// ErrNoCGO is returned by all operations when the package was built without CGo.
var ErrNoCGO = errors.New("ryugraph: built without CGo support")

// DB is a stub database handle used when CGo is unavailable.
type DB struct{}

// Conn is a stub connection handle used when CGo is unavailable.
type Conn struct{}

// Result is a stub query result used when CGo is unavailable.
type Result struct{}

func Open(path string) (*DB, error)                                          { return nil, ErrNoCGO }
func OpenMemory() (*DB, error)                                               { return nil, ErrNoCGO }
func OpenWithLimit(path string, bufferPoolBytes, maxDBBytes uint64) (*DB, error) {
	return nil, ErrNoCGO
}

func (db *DB) Close()                          {}
func (db *DB) Connect() (*Conn, error)         { return nil, ErrNoCGO }
func (c *Conn) Close()                         {}
func (c *Conn) Exec(cypher string) error       { return ErrNoCGO }
func (c *Conn) Query(cypher string) (*Result, error) { return nil, ErrNoCGO }
func (r *Result) Close()                       {}
func (r *Result) NumTuples() uint64            { return 0 }
func (r *Result) NumColumns() uint64           { return 0 }
func (r *Result) HasNext() bool                { return false }
func (r *Result) Next() []string               { return nil }

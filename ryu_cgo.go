//go:build cgo

// Package ryugraph provides a CGo wrapper around the RyuGraph embedded graph
// database (libryu). The wrapper exposes only the primitives needed for the
// sessiondata behavior graph: open/close a DB, execute Cypher, iterate results.
package ryugraph

/*
#cgo darwin  CFLAGS:  -I${SRCDIR}
#cgo darwin  LDFLAGS: -L${SRCDIR} -lryu_static -lc++
#cgo linux   CFLAGS:  -I${SRCDIR}
#cgo linux   LDFLAGS: -L${SRCDIR} -lryu_static -lstdc++ -lpthread -ldl -lm
#cgo windows CFLAGS:  -I${SRCDIR}
#cgo windows LDFLAGS: -L${SRCDIR} -lryu_static -lstdc++ -lws2_32

// Declare only the opaque struct types CGo needs — avoids the self-referencing
// typedef aliases in ryu.h that confuse the CGo type parser.
#include <stdint.h>
#include <stdbool.h>

typedef struct { void* _database; }               ryu_database;
typedef struct { void* _connection; }              ryu_connection;
typedef struct { void* _query_result; bool _owned; } ryu_query_result;
typedef struct { void* _flat_tuple;  bool _owned; } ryu_flat_tuple;
typedef struct { void* _value;       bool _owned; } ryu_value;

// Forward-declare the wrapper functions from ryu_wrap.c
int      wrap_ryu_database_init(const char* path, ryu_database* db);
int      wrap_ryu_database_init_nowal(const char* path, uint64_t buffer_pool_bytes, uint64_t max_db_bytes, ryu_database* db);
int      wrap_ryu_database_init_with_limit(const char* path, uint64_t buffer_pool_bytes, uint64_t max_db_bytes, ryu_database* db);
void     wrap_ryu_database_destroy(ryu_database* db);
int      wrap_ryu_connection_init(ryu_database* db, ryu_connection* conn);
void     wrap_ryu_connection_destroy(ryu_connection* conn);
int      wrap_ryu_connection_query(ryu_connection* conn, const char* q, ryu_query_result* res);
void     wrap_ryu_query_result_destroy(ryu_query_result* res);
int      wrap_ryu_query_result_is_success(ryu_query_result* res);
char*    wrap_ryu_query_result_get_error_message(ryu_query_result* res);
int      wrap_ryu_query_result_has_next(ryu_query_result* res);
int      wrap_ryu_query_result_get_next(ryu_query_result* res, ryu_flat_tuple* t);
uint64_t wrap_ryu_query_result_get_num_columns(ryu_query_result* res);
uint64_t wrap_ryu_query_result_get_num_tuples(ryu_query_result* res);
void     wrap_ryu_flat_tuple_destroy(ryu_flat_tuple* t);
int      wrap_ryu_flat_tuple_get_value(ryu_flat_tuple* t, uint64_t idx, ryu_value* v);
void     wrap_ryu_value_destroy(ryu_value* v);
char*    wrap_ryu_value_to_string(ryu_value* v);

#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// DB wraps a ryu_database handle.
type DB struct{ h C.ryu_database }

// Conn wraps a ryu_connection handle.
type Conn struct{ h C.ryu_connection }

// Result wraps a ryu_query_result handle.
type Result struct{ h C.ryu_query_result }

// Open creates or opens a RyuGraph database at path.
// Pass "" or ":memory:" for a fully in-memory database (no files on disk).
func Open(path string) (*DB, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var db DB
	if C.wrap_ryu_database_init(cpath, &db.h) != 0 {
		return nil, errors.New("ryugraph: failed to open database at " + path)
	}
	return &db, nil
}

// OpenMemory creates a fully in-memory RyuGraph database.
// Data exists only for the lifetime of the process — nothing written to disk.
func OpenMemory() (*DB, error) { return Open(":memory:") }

// OpenWithLimit opens (or creates) a database at path with explicit memory and
// disk limits.
//
// bufferPoolBytes — max RAM used as page cache. Hot pages stay in memory; cold
// pages are evicted to disk (for an on-disk DB) or cause an error (":memory:").
// Use a real directory path to get graceful disk spill instead of errors.
//
// maxDBBytes — max on-disk database size (0 = ryugraph default ~8TB virtual).
// Useful to bound disk growth when using a real path.
func OpenWithLimit(path string, bufferPoolBytes, maxDBBytes uint64) (*DB, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var db DB
	if C.wrap_ryu_database_init_nowal(cpath, C.uint64_t(bufferPoolBytes), C.uint64_t(maxDBBytes), &db.h) != 0 {
		return nil, errors.New("ryugraph: failed to open database at " + path)
	}
	return &db, nil
}

// Close releases the database.
func (db *DB) Close() { C.wrap_ryu_database_destroy(&db.h) }

// Connect opens a connection on db.
func (db *DB) Connect() (*Conn, error) {
	var conn Conn
	if C.wrap_ryu_connection_init(&db.h, &conn.h) != 0 {
		return nil, errors.New("ryugraph: failed to create connection")
	}
	return &conn, nil
}

// Close releases the connection.
func (c *Conn) Close() { C.wrap_ryu_connection_destroy(&c.h) }

// Exec runs a Cypher statement that returns no rows (DDL, CREATE, MERGE).
func (c *Conn) Exec(cypher string) error {
	r, err := c.Query(cypher)
	if err != nil {
		return err
	}
	r.Close()
	return nil
}

// Query runs a Cypher query and returns the result set.
func (c *Conn) Query(cypher string) (*Result, error) {
	cq := C.CString(cypher)
	defer C.free(unsafe.Pointer(cq))
	var res Result
	if C.wrap_ryu_connection_query(&c.h, cq, &res.h) != 0 {
		// Result may still be populated with an error message even on hard failure.
		if res.h._query_result != nil {
			msg := C.GoString(C.wrap_ryu_query_result_get_error_message(&res.h))
			C.wrap_ryu_query_result_destroy(&res.h)
			return nil, errors.New("ryugraph: " + msg)
		}
		return nil, errors.New("ryugraph: query failed: " + cypher)
	}
	if C.wrap_ryu_query_result_is_success(&res.h) == 0 {
		msg := C.GoString(C.wrap_ryu_query_result_get_error_message(&res.h))
		C.wrap_ryu_query_result_destroy(&res.h)
		return nil, errors.New("ryugraph: " + msg)
	}
	return &res, nil
}

// Close releases the result.
func (r *Result) Close() { C.wrap_ryu_query_result_destroy(&r.h) }

// NumTuples returns total row count.
func (r *Result) NumTuples() uint64 {
	return uint64(C.wrap_ryu_query_result_get_num_tuples(&r.h))
}

// NumColumns returns column count.
func (r *Result) NumColumns() uint64 {
	return uint64(C.wrap_ryu_query_result_get_num_columns(&r.h))
}

// HasNext returns true if more rows remain.
func (r *Result) HasNext() bool {
	return C.wrap_ryu_query_result_has_next(&r.h) != 0
}

// Next advances to the next row and returns column values as strings.
// Returns nil when exhausted.
func (r *Result) Next() []string {
	if !r.HasNext() {
		return nil
	}
	var tuple C.ryu_flat_tuple
	if C.wrap_ryu_query_result_get_next(&r.h, &tuple) != 0 {
		return nil
	}
	defer C.wrap_ryu_flat_tuple_destroy(&tuple)

	n := int(C.wrap_ryu_query_result_get_num_columns(&r.h))
	row := make([]string, n)
	for i := 0; i < n; i++ {
		var val C.ryu_value
		if C.wrap_ryu_flat_tuple_get_value(&tuple, C.uint64_t(i), &val) != 0 {
			continue
		}
		cstr := C.wrap_ryu_value_to_string(&val)
		row[i] = C.GoString(cstr)
		C.free(unsafe.Pointer(cstr))
		C.wrap_ryu_value_destroy(&val)
	}
	return row
}

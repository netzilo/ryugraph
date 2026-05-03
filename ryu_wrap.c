// ryu_wrap.c — thin C wrappers around ryu.h so CGo never sees the
// self-referencing typedef aliases that confuse the CGo type parser.
#include "ryu.h"
#include <stdlib.h>

ryu_system_config wrap_ryu_default_system_config() {
    return ryu_default_system_config();
}

int wrap_ryu_database_init(const char* path, ryu_database* db) {
    return (int)ryu_database_init(path, ryu_default_system_config(), db);
}

int wrap_ryu_database_init_with_config(const char* path, uint64_t buffer_pool_bytes, uint64_t max_db_bytes, bool no_wal, ryu_database* db) {
    ryu_system_config cfg = ryu_default_system_config();
    cfg.buffer_pool_size = buffer_pool_bytes;
    if (max_db_bytes > 0) {
        cfg.max_db_size = max_db_bytes;
    }
    cfg.no_wal = no_wal;
    return (int)ryu_database_init(path, cfg, db);
}

int wrap_ryu_database_init_nowal(const char* path, uint64_t buffer_pool_bytes, uint64_t max_db_bytes, ryu_database* db) {
    return wrap_ryu_database_init_with_config(path, buffer_pool_bytes, max_db_bytes, true, db);
}

int wrap_ryu_database_init_with_limit(const char* path, uint64_t buffer_pool_bytes, uint64_t max_db_bytes, ryu_database* db) {
    return wrap_ryu_database_init_with_config(path, buffer_pool_bytes, max_db_bytes, false, db);
}

void wrap_ryu_database_destroy(ryu_database* db) {
    ryu_database_destroy(db);
}

int wrap_ryu_connection_init(ryu_database* db, ryu_connection* conn) {
    return (int)ryu_connection_init(db, conn);
}

void wrap_ryu_connection_destroy(ryu_connection* conn) {
    ryu_connection_destroy(conn);
}

int wrap_ryu_connection_query(ryu_connection* conn, const char* q, ryu_query_result* res) {
    return (int)ryu_connection_query(conn, q, res);
}

void wrap_ryu_query_result_destroy(ryu_query_result* res) {
    ryu_query_result_destroy(res);
}

int wrap_ryu_query_result_is_success(ryu_query_result* res) {
    return ryu_query_result_is_success(res) ? 1 : 0;
}

char* wrap_ryu_query_result_get_error_message(ryu_query_result* res) {
    return ryu_query_result_get_error_message(res);
}

int wrap_ryu_query_result_has_next(ryu_query_result* res) {
    return ryu_query_result_has_next(res) ? 1 : 0;
}

int wrap_ryu_query_result_get_next(ryu_query_result* res, ryu_flat_tuple* t) {
    return (int)ryu_query_result_get_next(res, t);
}

uint64_t wrap_ryu_query_result_get_num_columns(ryu_query_result* res) {
    return ryu_query_result_get_num_columns(res);
}

uint64_t wrap_ryu_query_result_get_num_tuples(ryu_query_result* res) {
    return ryu_query_result_get_num_tuples(res);
}

void wrap_ryu_flat_tuple_destroy(ryu_flat_tuple* t) {
    ryu_flat_tuple_destroy(t);
}

int wrap_ryu_flat_tuple_get_value(ryu_flat_tuple* t, uint64_t idx, ryu_value* v) {
    return (int)ryu_flat_tuple_get_value(t, idx, v);
}

void wrap_ryu_value_destroy(ryu_value* v) {
    ryu_value_destroy(v);
}

char* wrap_ryu_value_to_string(ryu_value* v) {
    return ryu_value_to_string(v);
}

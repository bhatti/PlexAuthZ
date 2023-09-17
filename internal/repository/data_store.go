package repository

import "time"

// DataStore interface
type DataStore interface {
	// CreateTable creates underlying table if needed.
	CreateTable(
		baseTableName string,
		baseTableSuffix string,
	) (err error)

	// Size of records for tenant.
	Size(
		baseTableName string,
		baseTableSuffix string,
		tenant string,
		namespace string,
	) (size int64, err error)

	// Get finds records by ids.
	Get(
		baseTableName string,
		baseTableSuffix string,
		tenant string,
		namespace string,
		ids ...string,
	) (res map[string][]byte, err error)

	// Query searches records by predicates.
	Query(
		baseTableName string,
		baseTableSuffix string,
		tenant string,
		namespace string,
		predicate map[string]string,
		lastEvaluatedKeyStr string,
		limit int64,
	) (res map[string][]byte, nextKeyStr string, err error)

	// Create adds a new record.
	Create(
		baseTableName string,
		baseTableSuffix string,
		tenant string,
		namespace string,
		id string,
		value []byte,
		expiration time.Duration) (err error)

	// Update changes existing record.
	Update(
		baseTableName string,
		baseTableSuffix string,
		tenant string,
		namespace string,
		id string,
		version int64,
		value []byte,
		expiration time.Duration) (err error)

	// Delete removes existing record.
	Delete(
		baseTableName string,
		baseTableSuffix string,
		tenant string,
		namespace string,
		id string,
	) (err error)

	// ClearTable removes all records for tenant
	ClearTable(
		baseTableName string,
		baseTableSuffix string,
		tenant string,
		namespace string,
	) (err error)
}

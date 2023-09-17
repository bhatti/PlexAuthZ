package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	log "github.com/sirupsen/logrus"
	"time"
)

// BaseRepository - provides default persistence
type BaseRepository[T any] struct {
	baseTableName   string
	baseTableSuffix string
	builder         domain.Factory[T]
	expiration      time.Duration
	store           DataStore
}

// NewBaseRepository creates base repository
func NewBaseRepository[T any](
	store DataStore,
	baseTableName string,
	baseTableSuffix string,
	expiration time.Duration,
	builder domain.Factory[T],
) (*BaseRepository[T], error) {
	if err := store.CreateTable(
		baseTableName,
		baseTableSuffix); err != nil {
		return nil, err
	}
	return &BaseRepository[T]{
		store:           store,
		baseTableName:   baseTableName,
		baseTableSuffix: baseTableSuffix,
		builder:         builder,
		expiration:      expiration,
	}, nil
}

// GetByIDs retrieves objects by ids
func (r *BaseRepository[T]) GetByIDs(
	_ context.Context,
	organizationID string,
	namespace string,
	ids ...string) (res map[string]*T, err error) {
	bytesMap, err := r.store.Get(
		r.baseTableName,
		r.baseTableSuffix,
		organizationID,
		namespace,
		ids...)
	if err != nil {
		return nil, domain.NewNotFoundError(
			fmt.Sprintf("failed to get objects [%s %s %s %v] due to %s",
				r.baseTableName, organizationID, namespace, ids, err))
	}
	res = make(map[string]*T)
	for k, v := range bytesMap {
		obj := r.builder()
		err = json.Unmarshal(v, obj)
		if err != nil {
			return nil, domain.NewMarshalError(
				fmt.Sprintf("failed to unmarshal object %s [%s %s %s %v] due to %s", k,
					r.baseTableName, organizationID, namespace, ids, err))
		}
		res[k] = obj
	}
	return
}

// GetByID - finds an object by id
func (r *BaseRepository[T]) GetByID(
	ctx context.Context,
	organizationID string,
	namespace string,
	id string,
) (*T, error) {
	res, err := r.GetByIDs(ctx, organizationID, namespace, id)
	if err != nil {
		return nil, domain.NewNotFoundError(
			fmt.Sprintf("failed to find object [%s/%s/%s] due to %s",
				organizationID, namespace, id, err))
	}
	if res[id] == nil {
		return nil, domain.NewNotFoundError(
			fmt.Sprintf("object not found with [%s %s %s]",
				organizationID, namespace, id))
	}
	return res[id], nil
}

// Create adds a new item in the database.
func (r *BaseRepository[T]) Create(
	_ context.Context,
	organizationID string,
	namespace string,
	id string,
	obj *T,
	expiration time.Duration,
) (err error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return domain.NewMarshalError(
			fmt.Sprintf("failed to marshal object [%s %s %s %s] due to %s",
				r.baseTableName, organizationID, namespace, id, err))
	}
	log.WithFields(log.Fields{
		"Component":    "DDBBaseRepository",
		"Organization": organizationID,
		"Namespace":    namespace,
		"Id":           id,
		"Table":        r.baseTableName,
		"ObjSize":      len(b),
	}).
		Debugf("creating object")
	if expiration.Seconds() <= 0 {
		expiration = r.expiration
	}
	return r.store.Create(
		r.baseTableName,
		r.baseTableSuffix,
		organizationID,
		namespace,
		id,
		b,
		expiration)
}

// Update changes existing record in the database.
func (r *BaseRepository[T]) Update(
	_ context.Context,
	organizationID string,
	namespace string,
	id string,
	version int64,
	obj *T,
	expiration time.Duration,
) (err error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return domain.NewMarshalError(
			fmt.Sprintf("failed to marshal object [%s %s %s %s] due to %s",
				r.baseTableName, organizationID, namespace, id, err))
	}
	log.WithFields(log.Fields{
		"Component":    "DDBBaseRepository",
		"Organization": organizationID,
		"Namespace":    namespace,
		"Id":           id,
		"Table":        r.baseTableName,
		"ObjSize":      len(b),
	}).
		Debugf("updating object")
	if expiration.Seconds() <= 0 {
		expiration = r.expiration
	}
	return r.store.Update(
		r.baseTableName,
		r.baseTableSuffix,
		organizationID,
		namespace,
		id,
		version,
		b,
		expiration)
}

// Query - queries data
func (r *BaseRepository[T]) Query(
	_ context.Context,
	organizationID string,
	namespace string,
	predicate map[string]string,
	lastOffsetToken string,
	limit int64) (res []*T, nextOffsetToken string, err error) {
	matched, nextOffsetToken, err := r.store.Query(
		r.baseTableName,
		r.baseTableSuffix,
		organizationID,
		namespace,
		predicate,
		lastOffsetToken,
		limit)
	if err != nil {
		return res, "", domain.NewDatabaseError(
			fmt.Sprintf("failed to query objects [%s %s %s] due to %s",
				r.baseTableName, organizationID, namespace, err))
	}
	for _, b := range matched {
		obj := r.builder()
		err = json.Unmarshal(b, obj)
		if err != nil {
			return res, "", domain.NewMarshalError(
				fmt.Sprintf("failed to unmarshal object [%s %s %s] after query due to %s",
					r.baseTableName, organizationID, namespace, err))
		}
		res = append(res, obj)
	}
	return res, nextOffsetToken, nil
}

func (r *BaseRepository[T]) Delete(
	_ context.Context,
	organizationID string,
	namespace string,
	id string) error {
	log.WithFields(log.Fields{
		"Component":    "DDBBaseRepository",
		"Organization": organizationID,
		"Namespace":    namespace,
		"Id":           id,
	}).
		Debugf("deleting object")
	return r.store.Delete(
		r.baseTableName,
		r.baseTableSuffix,
		organizationID,
		namespace,
		id)
}

// Size of table
func (r *BaseRepository[T]) Size(
	_ context.Context,
	organizationID string,
	namespace string,
) (int64, error) {
	return r.store.Size(
		r.baseTableName,
		r.baseTableSuffix,
		organizationID,
		namespace,
	)
}

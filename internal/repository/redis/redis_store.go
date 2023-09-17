package redis

import (
	"fmt"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

// Store cache service
type Store struct {
	pool *redis.Pool
}

// NewRedisStore constructor for Redis store.
func NewRedisStore(
	config *domain.Config,
) (*Store, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	hostPort := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
	pool := &redis.Pool{
		MaxIdle:   config.Redis.PoolSize,
		MaxActive: config.Redis.PoolSize,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", hostPort)
			if err != nil {
				return nil, err
			}
			if config.Redis.Password != "" {
				if _, err := conn.Do("AUTH", config.Redis.Password); err != nil {
					_ = conn.Close()
					return nil, err
				}
			}
			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
	logrus.WithFields(
		logrus.Fields{
			"Component": "RedisStore",
			"Host":      config.Redis.Host,
			"Port":      config.Redis.Port,
		}).Debugf("connected to Redis")
	return &Store{
		pool: pool,
	}, nil
}

// CreateTable no-op function.
func (r *Store) CreateTable(
	_ string, // base-table
	_ string, // base suffix
) (err error) {
	// no-op
	return nil
}

// Size returns number of rows for tenant in table.
func (r *Store) Size(
	baseTableName string,
	baseTableSuffix string,
	tenant string,
	namespace string,
) (size int64, err error) {
	conn := r.pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	tableName := toTableName(baseTableName, baseTableSuffix, tenant, namespace)
	i, err := conn.Do("HLEN", tableName)
	if err != nil {
		return 0, err
	}
	return utils.ToInt64(i), nil
}

// Get finds record by ids.
func (r *Store) Get(
	baseTableName string,
	baseTableSuffix string,
	tenant string,
	namespace string,
	ids ...string,
) (res map[string][]byte, err error) {
	conn := r.pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	tableName := toTableName(baseTableName, baseTableSuffix, tenant, namespace)
	res = make(map[string][]byte)
	arr, err := toArray(conn.Do("HMGET", insert(tableName, ids)...))
	if err != nil {
		return nil, err
	}
	for i, a := range arr {
		res[ids[i]], err = redis.Bytes(a, err)
		if err != nil {
			return nil, err
		}
	}
	return
}

// Query queries records for given predicates.
func (r *Store) Query(
	baseTableName string,
	baseTableSuffix string,
	tenant string,
	namespace string,
	predicate map[string]string,
	offsetStr string,
	limit int64,
) (res map[string][]byte, nextOffset string, err error) {
	res = make(map[string][]byte)
	conn := r.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	tableName := toTableName(baseTableName, baseTableSuffix, tenant, namespace)
	arr, err := toArray(conn.Do("HGETALL", tableName))
	offset, _ := strconv.Atoi(offsetStr)
	if err != nil {
		return nil, "", err
	}
	for i := offset * 2; i < len(arr); i += 2 {
		name, err := redis.String(arr[i], nil)
		if err != nil {
			return nil, "", err
		}
		value, err := redis.Bytes(arr[i+1], nil)
		if err != nil {
			return nil, "", err
		}
		if !utils.MatchPredicate(value, predicate) {
			continue
		}
		res[name] = value
		if limit > 0 && len(res) >= int(limit) {
			break
		}
	}
	nextOffset = fmt.Sprintf("%d", max(limit, len(res)))
	return
}

// Create updates item in Redis table -- no difference between Create and Update.
func (r *Store) Create(
	baseTableName string,
	baseTableSuffix string,
	tenant string,
	namespace string,
	id string,
	value []byte,
	expiration time.Duration) (err error) {
	return r.Update(baseTableName, baseTableSuffix, tenant, namespace, id, 0, value, expiration)
}

// Update updates item in Redis table.
func (r *Store) Update(
	baseTableName string,
	baseTableSuffix string,
	tenant string,
	namespace string,
	id string,
	_ int64, // no version support
	value []byte,
	expiration time.Duration) (err error) {
	conn := r.pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	tableName := toTableName(baseTableName, baseTableSuffix, tenant, namespace)
	_, err = conn.Do("HSET", tableName, id, value)
	if err == nil && expiration.Seconds() > 0 {
		_, err = conn.Do("EXPIRE", tableName, expiration.Seconds())
	}
	return
}

// Delete removes cache entry
func (r *Store) Delete(
	baseTableName string,
	baseTableSuffix string,
	tenant string,
	namespace string,
	id string) (err error) {
	conn := r.pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	tableName := toTableName(baseTableName, baseTableSuffix, tenant, namespace)
	_, err = conn.Do("HDEL", tableName, id)
	return
}

// ClearTable removes all entries in table
func (r *Store) ClearTable(
	baseTableName string,
	baseTableSuffix string,
	tenant string,
	namespace string,
) (err error) {
	conn := r.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	tableName := toTableName(baseTableName, baseTableSuffix, tenant, namespace)
	arr, err := toArray(conn.Do("HGETALL", tableName))
	if err != nil {
		return err
	}
	for i := 0; i < len(arr); i += 2 {
		name, err := redis.String(arr[i], nil)
		if err != nil {
			return err
		}
		_, err = conn.Do("HDEL", tableName, name)
		if err != nil {
			return err
		}
	}
	logrus.WithFields(logrus.Fields{
		"Component": "RedisStore",
		"Table":     tableName,
		"Deleted":   len(arr),
	}).
		Debugf("deleting all objects in table")
	return
}

func toArray(i interface{}, err error) ([]interface{}, error) {
	if err != nil {
		return nil, err
	}
	switch arr := i.(type) {
	case []interface{}:
		return arr, nil
	}
	return nil, domain.NewValidationError(
		fmt.Sprintf("toArray found unexpected type for %v", i))
}

func insert(tableName string, ids []string) (res []interface{}) {
	res = append(res, tableName)
	for _, id := range ids {
		res = append(res, id)
	}
	return
}

func toTableName(
	baseTableName string,
	baseTableSuffix string,
	organizationID string,
	namespace string,
) string {
	if baseTableName == "Organization" {
		return baseTableName
	} else if namespace == "" {
		return fmt.Sprintf("%s__%s", baseTableName, organizationID)
	}
	return fmt.Sprintf("%s__%s__%s__%s", baseTableName, organizationID, namespace, baseTableSuffix)
}

func max(i int64, j int) int {
	ii := int(i)
	if ii > j {
		return ii
	}
	return j
}

package ddb

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/bhatti/PlexAuthZ/internal/domain"
	"github.com/bhatti/PlexAuthZ/internal/utils"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"time"
)

// Store cache service
type Store struct {
	config *domain.DynamoDBConfig
	ddbSvc *dynamodb.DynamoDB
}

// NewDDBStore constructor for DynamoDB table.
func NewDDBStore(
	config *domain.Config,
) (*Store, error) {
	// Initialize a session from environment variables
	if err := config.DynamoDB.Validate(); err != nil {
		return nil, err
	}
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(config.DynamoDB.AWSRegion), // Change this to your region
		Endpoint: aws.String(config.DynamoDB.Endpoint),
	})
	if err != nil {
		return nil, err
	}
	// Create DynamoDB client
	svc := dynamodb.New(sess)
	logrus.WithFields(
		logrus.Fields{
			"Component": "DDBStore",
		}).Debugf("connected to DDB")
	return &Store{config: &config.DynamoDB, ddbSvc: svc}, nil
}

// CreateTable helper method creates DDB table.
func (r *Store) CreateTable(
	baseTableName string,
	_ string, // suffix
) (err error) {
	if b, _ := r.tableExists(baseTableName); b {
		return nil
	}
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(baseTableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(r.config.TenantPartitionName),
				KeyType:       aws.String("HASH"), // Partition Key
			},
			{
				AttributeName: aws.String(r.config.IDName),
				KeyType:       aws.String("RANGE"), // Sort Key
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(r.config.TenantPartitionName),
				AttributeType: aws.String("S"), // S for string
			},
			{
				AttributeName: aws.String(r.config.IDName),
				AttributeType: aws.String("S"), // S for string
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(r.config.ReadCapacityUnits),
			WriteCapacityUnits: aws.Int64(r.config.WriteCapacityUnits),
		},
	}
	_, err = r.ddbSvc.CreateTable(input)
	if err != nil {
		return err
	}
	return r.ddbSvc.WaitUntilTableExists(&dynamodb.DescribeTableInput{
		TableName: aws.String(baseTableName),
	})
}

// Size returns number of items in table.
func (r *Store) Size(
	tableName string,
	baseTableSuffix string,
	baseTenant string,
	namespace string,
) (size int64, err error) {
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		Select:    aws.String("COUNT"),
		KeyConditions: map[string]*dynamodb.Condition{
			r.config.TenantPartitionName: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(toTenant(baseTableSuffix, baseTenant, namespace)),
					},
				},
			},
		},
	}

	result, err := r.ddbSvc.Query(input)
	if err != nil {
		return 0, err
	}
	return *result.Count, nil
}

// Get finds items by ids.
func (r *Store) Get(
	baseTableName string,
	baseTableSuffix string,
	baseTenant string,
	namespace string,
	ids ...string,
) (res map[string][]byte, err error) {
	keys := make([]map[string]*dynamodb.AttributeValue, len(ids))
	res = make(map[string][]byte)

	for i, id := range ids {
		keys[i] = map[string]*dynamodb.AttributeValue{
			r.config.TenantPartitionName: {
				S: aws.String(toTenant(baseTableSuffix, baseTenant, namespace)),
			},
			r.config.IDName: {
				S: aws.String(id),
			},
		}
	}
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			baseTableName: {
				Keys: keys,
			},
		},
	}

	result, err := r.ddbSvc.BatchGetItem(input)
	if err != nil {
		return nil, err
	}
	for _, item := range result.Responses[baseTableName] {
		id, val := r.getItem(item)
		if val != nil {
			res[id] = val
		}
	}
	for _, id := range ids {
		if res[id] == nil {
			return nil, domain.NewNotFoundError(
				fmt.Sprintf("failed to get object for id %s in %s", id, baseTableName))
		}
	}
	return res, nil
}

// Query searches items by predicates - not very optimal -- better use GSI in the future.
func (r *Store) Query(
	baseTableName string,
	baseTableSuffix string,
	baseTenant string,
	namespace string,
	predicate map[string]string,
	lastEvaluatedKeyStr string,
	limit int64,
) (res map[string][]byte, nextKeyStr string, err error) {
	if limit == 0 {
		limit = 500
	}
	res = make(map[string][]byte)
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue
	if lastEvaluatedKeyStr != "" {
		decoded, err := base64.StdEncoding.DecodeString(lastEvaluatedKeyStr)
		if err == nil {
			err = json.Unmarshal(decoded, &lastEvaluatedKey)
		}
		if err != nil {
			log.WithFields(log.Fields{
				"Error":               err,
				"LastEvaluatedKeyStr": lastEvaluatedKeyStr,
			}).Warnf("failed to parse lastEvaluatedKeyStr")
		}
	}

	input := &dynamodb.QueryInput{
		TableName: aws.String(baseTableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"tenant": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(toTenant(baseTableSuffix, baseTenant, namespace)),
					},
				},
			},
		},
		ExclusiveStartKey: lastEvaluatedKey,
		Limit:             aws.Int64(limit),
	}

	result, err := r.ddbSvc.Query(input)
	if err != nil {
		return nil, "", err
	}

	for _, item := range result.Items {
		id, val := r.getItem(item)
		if val != nil {
			if utils.MatchPredicate(val, predicate) {
				res[id] = val
			}
		}
	}

	if result.LastEvaluatedKey != nil {
		jsonKey, err := json.Marshal(result.LastEvaluatedKey)
		if err != nil {
			return nil, "", err
		}
		nextKeyStr = base64.StdEncoding.EncodeToString(jsonKey)
	}

	return
}

// Create adds a new item.
func (r *Store) Create(
	baseTableName string,
	baseTableSuffix string,
	baseTenant string,
	namespace string,
	id string,
	value []byte,
	expiration time.Duration) (err error) {
	expireTime := time.Now().Add(expiration).Unix()
	// new insert
	input := &dynamodb.PutItemInput{
		TableName: aws.String(baseTableName),
		Item: map[string]*dynamodb.AttributeValue{
			r.config.TenantPartitionName: {
				S: aws.String(toTenant(baseTableSuffix, baseTenant, namespace)),
			},
			r.config.IDName: {
				S: aws.String(id),
			},
			"val": {
				S: aws.String(string(value)),
			},
			"ver": {
				N: aws.String("1"),
			},
		},
		ConditionExpression: aws.String("attribute_not_exists(tenant) AND attribute_not_exists(id)"),
	}
	if expiration.Seconds() > 0 {
		input.Item["expire_at"] = &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", expireTime))}
	}
	_, err = r.ddbSvc.PutItem(input) // check for ConditionalCheckFailedException
	return
}

// Update updates existing entry.
func (r *Store) Update(
	baseTableName string,
	baseTableSuffix string,
	baseTenant string,
	namespace string,
	id string,
	version int64,
	value []byte,
	expiration time.Duration) (err error) {
	expireTime := time.Now().Add(expiration).Unix()
	if version < 0 {
		// no-version assumed so just treat as writing either new or existing.
		input := &dynamodb.PutItemInput{
			TableName: aws.String(baseTableName),
			Item: map[string]*dynamodb.AttributeValue{
				r.config.TenantPartitionName: {
					S: aws.String(toTenant(baseTableSuffix, baseTenant, namespace)),
				},
				r.config.IDName: {
					S: aws.String(id),
				},
				"val": {
					S: aws.String(string(value)),
				},
				"ver": {
					N: aws.String("1"),
				},
			},
		}
		if expiration.Seconds() > 0 {
			input.Item["expire_at"] = &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", expireTime))}
		}
		_, err = r.ddbSvc.PutItem(input) // check for ConditionalCheckFailedException
	} else {
		// update existing record
		condition := "ver = :expectedVersion"
		expressionAttributeValues := map[string]*dynamodb.AttributeValue{
			":val": {
				S: aws.String(string(value)),
			},
			":ver": {
				N: aws.String(fmt.Sprintf("%d", version+1)),
			},
			":expectedVersion": {
				N: aws.String(fmt.Sprintf("%d", version)),
			},
		}

		input := &dynamodb.UpdateItemInput{
			TableName: aws.String(baseTableName),
			Key: map[string]*dynamodb.AttributeValue{
				r.config.TenantPartitionName: {S: aws.String(toTenant(baseTableSuffix, baseTenant, namespace))},
				r.config.IDName:              {S: aws.String(id)},
			},
			ConditionExpression:       aws.String(condition),
			UpdateExpression:          aws.String("SET val = :val, ver = :ver"),
			ExpressionAttributeValues: expressionAttributeValues,
		}
		if expiration.Seconds() > 0 {
			expressionAttributeValues["expire_at"] = &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", expireTime))}
		}

		_, err = r.ddbSvc.UpdateItem(input)
	}
	return
}

// Delete removes existing item in DDB table.
func (r *Store) Delete(
	baseTableName string,
	baseTableSuffix string,
	baseTenant string,
	namespace string,
	id string,
) (err error) {
	_, err = r.ddbSvc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(baseTableName),
		Key: map[string]*dynamodb.AttributeValue{
			r.config.TenantPartitionName: {S: aws.String(toTenant(baseTableSuffix, baseTenant, namespace))},
			r.config.IDName:              {S: aws.String(id)},
		},
	})
	return
}

// ClearTable removes all entries for tenant in table.
func (r *Store) ClearTable(
	baseTableName string,
	baseTableSuffix string,
	baseTenant string,
	namespace string,
) (err error) {
	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(baseTableName),
		KeyConditions: map[string]*dynamodb.Condition{
			r.config.TenantPartitionName: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(toTenant(baseTableSuffix, baseTenant, namespace)),
					},
				},
			},
		},
	}

	queryOut, err := r.ddbSvc.Query(queryInput)
	if err != nil {
		return err
	}

	for _, item := range queryOut.Items {
		deleteInput := &dynamodb.DeleteItemInput{
			TableName: aws.String(baseTableName),
			Key: map[string]*dynamodb.AttributeValue{
				r.config.TenantPartitionName: {
					S: aws.String(toTenant(baseTableSuffix, baseTenant, namespace)),
				},
				r.config.IDName: item[r.config.IDName],
			},
		}

		_, err := r.ddbSvc.DeleteItem(deleteInput)
		if err != nil {
			return err
		}
	}
	return
}

func (r *Store) getItem(
	item map[string]*dynamodb.AttributeValue,
) (string, []byte) {
	id := item[r.config.IDName]
	value := item["val"]
	if id != nil && id.S != nil &&
		value != nil && value.S != nil {
		return *id.S, []byte(*value.S)
	}
	return "", nil
}

func (r *Store) tableExists(
	tableName string,
) (bool, error) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	_, err := r.ddbSvc.DescribeTable(input)
	if err != nil {
		if err, ok := err.(awserr.Error); ok {
			switch err.Code() {
			case dynamodb.ErrCodeResourceNotFoundException:
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

func toTenant(
	baseTableSuffix string,
	organizationID string,
	namespace string,
) string {
	return fmt.Sprintf("%s__%s__%s", organizationID, namespace, baseTableSuffix)
}

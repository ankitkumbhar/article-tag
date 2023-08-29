package model

import (
	"article-tag/internal/constant"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// tableName
const tableName = "article-follow-tag-v5"

// dynamoAPI
type dynamoAPI interface {
	DescribeTable(ctx context.Context, params *dynamodb.DescribeTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
	CreateTable(ctx context.Context, params *dynamodb.CreateTableInput, optFns ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

type tag struct {
	db     dynamoAPI
	logger *zap.Logger
	// db dynamodb.Client
}

func NewTag(m dynamoAPI, logger *zap.Logger) UserTagStore {
	return &tag{db: m, logger: logger}
}

// DescribeTable
func (t *tag) DescribeTable(ctx context.Context) error {
	input := dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	_, err := t.db.DescribeTable(ctx, &input)
	if err != nil {
		return err
	}

	return nil
}

// CreateTable
func (t *tag) CreateTable(ctx context.Context) error {
	i := dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("PK"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("SK"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("TagCount"),
			AttributeType: types.ScalarAttributeTypeN,
		},
			// {
			// 	AttributeName: aws.String("TagID"),
			// 	AttributeType: types.ScalarAttributeTypeS,
			// },
			{
				AttributeName: aws.String("TagName"),
				AttributeType: types.ScalarAttributeTypeS,
			}, {
				AttributeName: aws.String("CreatedAt"),
				AttributeType: types.ScalarAttributeTypeS,
			}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("PK"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("SK"),
			KeyType:       types.KeyTypeRange,
		}},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{{
			IndexName: aws.String("TagIndex"),
			KeySchema: []types.KeySchemaElement{{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			}, {
				AttributeName: aws.String("TagCount"),
				KeyType:       types.KeyTypeRange,
			}},
			Projection: &types.Projection{
				ProjectionType:   types.ProjectionTypeInclude,
				NonKeyAttributes: []string{"TagID", "TagName"},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		}},
		LocalSecondaryIndexes: []types.LocalSecondaryIndex{
			{
				IndexName: aws.String("LSI1"),
				KeySchema: []types.KeySchemaElement{{
					AttributeName: aws.String("PK"),
					KeyType:       types.KeyTypeHash,
				}, {
					AttributeName: aws.String("TagName"),
					KeyType:       types.KeyTypeRange,
				}},
				Projection: &types.Projection{
					ProjectionType:   types.ProjectionTypeInclude,
					NonKeyAttributes: []string{"TagName"},
				},
			},
			{
				IndexName: aws.String("LSI2"),
				KeySchema: []types.KeySchemaElement{{
					AttributeName: aws.String("PK"),
					KeyType:       types.KeyTypeHash,
				}, {
					AttributeName: aws.String("CreatedAt"),
					KeyType:       types.KeyTypeRange,
				}},
				Projection: &types.Projection{
					ProjectionType:   types.ProjectionTypeInclude,
					NonKeyAttributes: []string{"CreatedAt"},
				},
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := t.db.CreateTable(ctx, &i)
	if err != nil {
		t.logger.Error("error creating table", zap.Error(err))
		return err
	}

	return nil
}

func (t *tag) Store(ctx context.Context, username, publication, tagName, tagID string) error {
	item := UserTag{
		PK:          fmt.Sprintf("%v#%v", username, publication),
		SK:          tagID,
		TagID:       tagID,
		TagName:     tagName,
		CreatedAt:   time.Now().UTC().Format(time.RFC3339Nano),
		Username:    username,
		Publication: publication,
	}

	// convert struct to map
	inputMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		t.logger.Error("marshal failed", zap.Error(err))
		return err
	}

	input2 := dynamodb.PutItemInput{
		TableName:    aws.String(tableName),
		Item:         inputMap,
		ReturnValues: types.ReturnValueAllOld, // used to get old content
	}

	putItemOutput, err := t.db.PutItem(ctx, &input2)
	if err != nil {
		return err
	}

	// update popular tag count only if the user is following new tag
	// means when user follows already followed tag we dont need to update count
	if putItemOutput.Attributes == nil {
		input3 := dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", item.Publication)},
				"SK": &types.AttributeValueMemberS{Value: item.TagID},
			},
			UpdateExpression: aws.String("SET TagCount = if_not_exists(TagCount, :v1) + :incr, TagID = :v2, TagName = :v3"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v1":   &types.AttributeValueMemberN{Value: "0"},
				":incr": &types.AttributeValueMemberN{Value: "1"},
				":v2":   &types.AttributeValueMemberS{Value: item.TagID},
				":v3":   &types.AttributeValueMemberS{Value: item.TagName},
			},
		}

		_, err = t.db.UpdateItem(ctx, &input3)
		if err != nil {
			t.logger.Error("error updating tag counter while storing item", zap.Error(err))
			return err
		}
	}

	return nil
}

func (t *tag) Get(ctx context.Context, username, publication, order string) ([]*UserTag, error) {
	// get indexname and scanIndex using order
	indexName, scanIndex := getIndexNameAndScanOrder(order)

	t.logger.Debug("selected indexName and scanOrder", zapcore.Field{Key: "index_name", Type: zapcore.StringType,
		String: indexName}, zapcore.Field{Key: "scan_index", Type: zapcore.BoolType, Interface: scanIndex})

	queryInput := dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		IndexName:              aws.String(indexName),
		KeyConditionExpression: aws.String("#v1 = :v1"),
		ExpressionAttributeNames: map[string]string{
			"#v1": "PK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s", username, publication)},
		},
		ScanIndexForward:     aws.Bool(scanIndex),
		ProjectionExpression: aws.String("PK, SK, TagID, TagName"),
	}

	// fetch item
	res, err := t.db.Query(ctx, &queryInput)
	if err != nil {
		return nil, err
	}

	userTags := []*UserTag{}
	for _, val := range res.Items {
		var m UserTag

		err := attributevalue.UnmarshalMap(val, &m)
		if err != nil {
			t.logger.Error("unmarshal failed while fetching user tags", zap.Error(err))
			return nil, err
		}

		userTags = append(userTags, &UserTag{
			TagID:   m.TagID,
			TagName: m.TagName,
		})
	}

	return userTags, nil
}

// getIndexNameAndScanOrder
func getIndexNameAndScanOrder(order string) (string, bool) {

	var (
		indexName = "LSI1"
		scanIndex = true
	)

	/*
		Order:
		default - [indexName:LSI1, scanIndex: true]
		createdatdesc - [indexName:LSI2, scanIndex:false]
		createdatasc - [indexName:LSI2, scanIndex:true]
		tagname - [indexname:LSI1, scanIndex:false]
	*/

	switch order {
	case constant.CreatedAtDesc:
		indexName = "LSI2"
		scanIndex = false

	case constant.CreatedAtAsc:
		indexName = "LSI2"

	case constant.TagName:
		scanIndex = false
	default:

		scanIndex = true
	}

	return indexName, scanIndex
}

func (t *tag) Delete(ctx context.Context, username, publication, tagID, tagName string) error {

	input := dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s", username, publication)},
			"SK": &types.AttributeValueMemberS{Value: tagID},
		},
		ConditionExpression: aws.String("#v1 = :v1"),
		ExpressionAttributeNames: map[string]string{
			"#v1": "TagName",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: tagName},
		},
		ReturnValues: types.ReturnValueAllOld,
	}

	// fetch item
	delItemResp, err := t.db.DeleteItem(ctx, &input)
	if err != nil {
		return err
	}

	// when delete item response contains attribute
	// decrement the popularity count of deleted tag
	if delItemResp.Attributes != nil {
		input3 := dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", publication)},
				"SK": &types.AttributeValueMemberS{Value: tagID},
			},
			UpdateExpression: aws.String("SET TagCount = TagCount - :decr"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":decr": &types.AttributeValueMemberN{Value: "1"},
			},
		}

		_, err = t.db.UpdateItem(ctx, &input3)
		if err != nil {
			t.logger.Error("error updating tag counter while deleting item", zap.Error(err))
			return err
		}
	}

	return nil
}

func (t *tag) GetPopularTags(ctx context.Context, username, publication string) ([]string, error) {
	// Steps:
	// 1. Fetch the existing tags of the user
	// 2. To get popular tags, excluded the existing tags using filterExpression
	var (
		existingTags []*UserTag
		err          error
		hasMoreItems = true
		userTags     = []string{}
	)

	// check for empty username, when no username passed
	// return all tags of that particular publication
	if username != "" {
		existingTags, err = t.Get(ctx, username, publication, "")
		if err != nil {
			return nil, err
		}
	}

	var exclusiveStartKey map[string]types.AttributeValue = nil

	// iterate until we fetch all items
	for hasMoreItems {
		queryInput := dynamodb.QueryInput{
			TableName:              aws.String(tableName),
			IndexName:              aws.String("TagIndex"),
			KeyConditionExpression: aws.String("#v1 = :v1 AND #v2 > :v2"),
			ExpressionAttributeNames: map[string]string{
				"#v1": "PK",
				"#v2": "TagCount",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v1": &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", publication)},
				":v2": &types.AttributeValueMemberN{Value: "0"},
			},
			ScanIndexForward: aws.Bool(false),
			Limit:            aws.Int32(constant.PopularTagLimit),
		}

		// once we received lastEvaluatedKey from previous iteration
		// add it as exclusive start key for next iteration
		if exclusiveStartKey != nil {

			sk := ExclusiveStartKey{}
			err = attributevalue.UnmarshalMap(exclusiveStartKey, &sk)
			if err != nil {
				t.logger.Error("unmarshal failed", zap.Error(err))
				return nil, err
			}

			// queryInput.ExclusiveStartKey = starKey
			queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
				"PK":       &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", publication)},
				"SK":       &types.AttributeValueMemberS{Value: sk.SK},
				"TagCount": &types.AttributeValueMemberN{Value: sk.TagCount},
			}

			// {"PK":{"S":"PUB#AK"}, "SK":{"S":"2"},"TagCount":{"N":"1"}}
		}

		// exclude existing tags from the popular tags
		if len(existingTags) > 0 {
			prepareFilterExpression(&queryInput, existingTags)
		}

		// fetch item
		res, err := t.db.Query(ctx, &queryInput)
		if err != nil {
			return nil, err
		}

		// store the last evaluated key
		// used to iterate to fetch remaining items
		exclusiveStartKey = res.LastEvaluatedKey

		// fetch the tags
		for _, val := range res.Items {
			var m UserTag

			err := attributevalue.UnmarshalMap(val, &m)
			if err != nil {
				t.logger.Error("unmarshal failed", zap.Error(err))
				return nil, err
			}

			userTags = append(userTags, m.TagName)
		}

		// break the loop once the last item is fetched
		if res.LastEvaluatedKey == nil {
			// set false to break the loop
			hasMoreItems = false
		}
	}

	return userTags, nil
}

func prepareFilterExpression(queryInput *dynamodb.QueryInput, existingTags []*UserTag) {
	// filterExpression := fmt.Sprintf("%s  NOT (Tag IN (", *queryInput.FilterExpression)
	filterExpression := fmt.Sprintf("NOT (TagName IN (")

	filterAttr := []string{}

	for k, val := range existingTags {

		// prepare key
		key := fmt.Sprintf(":exclude%v", k)

		filterAttr = append(filterAttr, key)

		queryInput.ExpressionAttributeValues[key] = &types.AttributeValueMemberS{Value: val.TagName}
	}

	// join the filter expression placeholder
	filterExpression += strings.Join(filterAttr, ", ")

	// end of bracket
	filterExpression += "))"

	// update the filter expression
	queryInput.FilterExpression = aws.String(filterExpression)
}

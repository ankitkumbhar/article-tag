package model

import (
	"article-tag/internal/constant"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// tableName
const tableName = "article-follow-tag"

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
	db dynamoAPI
	// db dynamodb.Client
}

func NewTag(m dynamoAPI) UserTagStore {
	return &tag{db: m}
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
			AttributeName: aws.String("TotalCount"),
			AttributeType: types.ScalarAttributeTypeN,
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
				AttributeName: aws.String("TotalCount"),
				KeyType:       types.KeyTypeRange,
			}},
			Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		}},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := t.db.CreateTable(ctx, &i)
	if err != nil {
		fmt.Println("Error creating table", err)
		return err
	}

	return nil
}

func (t *tag) Store(ctx context.Context, item *UserTag) error {
	// convert struct to map
	inputMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		log.Fatal("marshal failed", err)
		return err
	}

	input2 := dynamodb.PutItemInput{
		TableName:    aws.String(tableName),
		Item:         inputMap,
		ReturnValues: types.ReturnValueAllOld, // used to get old content
	}

	putItemOutput, err := t.db.PutItem(ctx, &input2)
	if err != nil {
		fmt.Println("error storing new item: ", err)
		return err
	}

	// update popular tag count only if the user is following new tag
	// means when user follows already followed tag we dont need to update count
	if putItemOutput.Attributes == nil {
		input3 := dynamodb.UpdateItemInput{
			TableName: aws.String(tableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", item.Publication)},
				"SK": &types.AttributeValueMemberS{Value: item.Tag},
			},
			UpdateExpression: aws.String("SET TotalCount = if_not_exists(TotalCount, :v1) + :incr, Tag = :v2"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v1":   &types.AttributeValueMemberN{Value: "0"},
				":incr": &types.AttributeValueMemberN{Value: "1"},
				":v2":   &types.AttributeValueMemberS{Value: item.Tag},
			},
		}

		_, err = t.db.UpdateItem(ctx, &input3)
		if err != nil {
			fmt.Println("error updating tag counter while storing item", err)
			return err
		}
	}

	return nil
}

func (t *tag) Get(ctx context.Context, item *UserTag) ([]string, error) {
	queryInput := dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("#v1 = :v1"),
		ExpressionAttributeNames: map[string]string{
			"#v1": "PK",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s", item.Username, item.Publication)},
		},
	}

	// fetch item
	res, err := t.db.Query(ctx, &queryInput)
	if err != nil {
		return nil, err
	}

	userTags := []string{}
	for _, val := range res.Items {
		var m UserTag

		err := attributevalue.UnmarshalMap(val, &m)
		if err != nil {
			log.Fatal("unmarshal failed", err)
		}

		userTags = append(userTags, m.Tag)
	}

	return userTags, nil
}

func (t *tag) Delete(ctx context.Context, item *UserTag) error {
	input := dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s", item.Username, item.Publication)},
			"SK": &types.AttributeValueMemberS{Value: item.Tag},
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
				"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", item.Publication)},
				"SK": &types.AttributeValueMemberS{Value: item.Tag},
			},
			UpdateExpression: aws.String("SET TotalCount = TotalCount - :decr"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":decr": &types.AttributeValueMemberN{Value: "1"},
			},
		}

		_, err = t.db.UpdateItem(ctx, &input3)
		if err != nil {
			fmt.Println("error updating tag counter while deleting item : ", err)
			return err
		}
	}

	return nil
}

func (t *tag) GetPopularTags(ctx context.Context, item *UserTag) ([]string, error) {
	// Steps:
	// 1. Fetch the existing tags of the user
	// 2. To get popular tags, excluded the existing tags using filterExpression
	var (
		existingTags []string
		err          error
		hasMoreItems = true
		userTags     = []string{}
	)

	// check for empty username, when no username passed
	// return all tags of that particular publication
	if item.Username != "" {
		existingTags, err = t.Get(ctx, item)
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
				"#v2": "TotalCount",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":v1": &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", item.Publication)},
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
				log.Fatal("unmarshal failed", err)
			}

			// queryInput.ExclusiveStartKey = starKey
			queryInput.ExclusiveStartKey = map[string]types.AttributeValue{
				"PK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", item.Publication)},
				"SK":         &types.AttributeValueMemberS{Value: sk.SK},
				"TotalCount": &types.AttributeValueMemberN{Value: sk.TotalCount},
			}
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
				log.Fatal("unmarshal failed", err)
			}

			userTags = append(userTags, m.SK)
		}

		// break the loop once the last item is fetched
		if res.LastEvaluatedKey == nil {
			// set false to break the loop
			hasMoreItems = false
		}
	}

	return userTags, nil
}

func prepareFilterExpression(queryInput *dynamodb.QueryInput, existingTags []string) {
	// filterExpression := fmt.Sprintf("%s  NOT (Tag IN (", *queryInput.FilterExpression)
	filterExpression := fmt.Sprintf("NOT (Tag IN (")

	filterAttr := []string{}

	for k, val := range existingTags {

		// prepare key
		key := fmt.Sprintf(":exclude%v", k)

		filterAttr = append(filterAttr, key)

		queryInput.ExpressionAttributeValues[key] = &types.AttributeValueMemberS{Value: val}
	}

	// join the filter expression placeholder
	filterExpression += strings.Join(filterAttr, ", ")

	// end of bracket
	filterExpression += "))"

	// update the filter expression
	queryInput.FilterExpression = aws.String(filterExpression)
}

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
		TableName: aws.String("article-follow-tag-2"),
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
		TableName: aws.String("article-follow-tag-2"),
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("PK#PUB"),
			AttributeType: types.ScalarAttributeTypeS,
		}, {
			AttributeName: aws.String("SK"),
			AttributeType: types.ScalarAttributeTypeS,
		},
			// {
			// 	AttributeName: aws.String("Publication"),
			// 	AttributeType: types.ScalarAttributeTypeS,
			// },
			// {
			// 	AttributeName: aws.String("Tag"),
			// 	AttributeType: types.ScalarAttributeTypeS,
			// },
			{
				AttributeName: aws.String("TotalCount"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("PK#PUB"),
			KeyType:       types.KeyTypeHash,
		}, {
			AttributeName: aws.String("SK"),
			KeyType:       types.KeyTypeRange,
		}},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{{
			IndexName: aws.String("TagIndex"),
			KeySchema: []types.KeySchemaElement{{
				AttributeName: aws.String("PK#PUB"),
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
	// inputMap, err := attributevalue.MarshalMap(item)
	// if err != nil {
	// 	log.Fatal("marshal failed", err)
	// 	return err
	// }

	// inputMap["PK#PUB"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%v#%v", item.Username, item.Publication)}
	// inputMap["SK"] = &types.AttributeValueMemberS{Value: item.Tag}

	input2 := dynamodb.PutItemInput{
		TableName: aws.String("article-follow-tag-2"),
		Item: map[string]types.AttributeValue{
			"PK#PUB":      &types.AttributeValueMemberS{Value: fmt.Sprintf("%v#%v", item.Username, item.Publication)},
			"SK":          &types.AttributeValueMemberS{Value: item.Tag},
			"Publication": &types.AttributeValueMemberS{Value: item.Publication},
			"Tag":         &types.AttributeValueMemberS{Value: item.Tag},
		},
		// Item:         inputMap,
		ReturnValues: types.ReturnValueAllOld, // used to get old content
	}

	putItemOutput, err := t.db.PutItem(ctx, &input2)
	if err != nil {
		fmt.Println("error storing new item 1111: ", err)
		return err
	}

	// update popular tag count only if the user is following new tag
	// means when user follows already followed tag we dont need to update count
	if putItemOutput.Attributes == nil {
		input3 := dynamodb.UpdateItemInput{
			TableName: aws.String("article-follow-tag-2"),
			Key: map[string]types.AttributeValue{
				"PK#PUB": &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", item.Publication)},
				"SK":     &types.AttributeValueMemberS{Value: item.Tag},
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
		TableName:              aws.String("article-follow-tag-2"),
		KeyConditionExpression: aws.String("#v1 = :v1"),
		ExpressionAttributeNames: map[string]string{
			"#v1": "PK#PUB",
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
		TableName: aws.String("article-follow-tag-2"),
		Key: map[string]types.AttributeValue{
			"PK#PUB": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s", item.Username, item.Publication)},
			"SK":     &types.AttributeValueMemberS{Value: item.Tag},
		},
		ReturnValues: types.ReturnValueAllOld,
	}

	// fetch item
	_, err := t.db.DeleteItem(ctx, &input)
	if err != nil {
		return err
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
	)

	// check for empty username, when no username passed
	// return all tags of that particular publication
	if item.Username != "" {
		existingTags, err = t.Get(ctx, item)
		if err != nil {
			return nil, err
		}
	}

	queryInput := dynamodb.QueryInput{
		TableName:              aws.String("article-follow-tag-2"),
		IndexName:              aws.String("TagIndex"),
		KeyConditionExpression: aws.String("#v1 = :v1"),
		// FilterExpression:       aws.String("NOT (Tag IN :v3)"),
		// FilterExpression: aws.String(""),
		// FilterExpression: aws.String("#v2 > :v2"),
		ExpressionAttributeNames: map[string]string{
			"#v1": "PK#PUB",
			// "#v2": "TotalCount",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", item.Publication)},
			// ":v2": &types.AttributeValueMemberN{Value: strconv.Itoa(constant.PopularTagCount)},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(constant.PopularTagLimit),
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

	userTags := []string{}
	for _, val := range res.Items {
		var m UserTag

		err := attributevalue.UnmarshalMap(val, &m)
		if err != nil {
			log.Fatal("unmarshal failed", err)
		}

		userTags = append(userTags, m.SK)

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

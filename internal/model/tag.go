package model

import (
	"article-tag/internal/constant"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type tag struct {
	db *dynamodb.Client
}

func (t *tag) Store(ctx context.Context, item *UserTag) error {
	input2 := dynamodb.PutItemInput{
		TableName: aws.String("article-tag-6"),
		Item: map[string]types.AttributeValue{
			"PK#PUB":      &types.AttributeValueMemberS{Value: fmt.Sprintf("%v#%v", item.Username, item.Publication)},
			"SK":          &types.AttributeValueMemberS{Value: item.Tag},
			"Publication": &types.AttributeValueMemberS{Value: item.Publication},
			"Tag":         &types.AttributeValueMemberS{Value: item.Tag},
			// "CreatedAt":   &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
		// Item: map[string]types.AttributeValue{
		// 	"PK#PUB":      &types.AttributeValueMemberS{Value: "Sandip#AK"},
		// 	"SK":          &types.AttributeValueMemberS{Value: "tag60"},
		// 	"Publication": &types.AttributeValueMemberS{Value: "AK"},
		// 	"Tag":         &types.AttributeValueMemberS{Value: "tag60"},
		// 	// "CreatedAt":   &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		// },
	}
	_, err := t.db.PutItem(ctx, &input2)
	if err != nil {
		fmt.Println("error storing new item 1111: ", err)
		return err
	}

	// input3 := dynamodb.PutItemInput{
	// 	TableName: aws.String("article-tag-6"),
	// 	Item: map[string]types.AttributeValue{
	// 		"PK#PUB": &types.AttributeValueMemberS{Value: "PUB#AK"},
	// 		"SK":     &types.AttributeValueMemberS{Value: "tag60"},
	// 		// "Tag":    &types.AttributeValueMemberS{Value: "tag60"},
	// 		"TotalCount": &types.AttributeValueMemberN{Value: "1"},
	// 	},
	// }

	input3 := dynamodb.UpdateItemInput{
		TableName: aws.String("article-tag-6"),
		Key: map[string]types.AttributeValue{
			"PK#PUB": &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", item.Publication)},
			"SK":     &types.AttributeValueMemberS{Value: item.Tag},
		},
		// Key: map[string]types.AttributeValue{
		// 	"PK#PUB": &types.AttributeValueMemberS{Value: "PUB#AK"},
		// 	"SK":     &types.AttributeValueMemberS{Value: "tag60"},
		// },
		UpdateExpression: aws.String("SET TotalCount = if_not_exists(TotalCount, :v1) + :incr, Tag = :v2"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1":   &types.AttributeValueMemberN{Value: "0"},
			":incr": &types.AttributeValueMemberN{Value: "1"},
			":v2":   &types.AttributeValueMemberS{Value: item.Tag},
		},
	}
	_, err = t.db.UpdateItem(ctx, &input3)
	// _, err = t.db.PutItem(ctx, &input3)
	if err != nil {
		fmt.Println("error storing new item 2222: ", err)
		return err
	}

	return nil
}

func (t *tag) Get(ctx context.Context, item *UserTag) ([]string, error) {
	queryInput := dynamodb.QueryInput{
		TableName:              aws.String("article-tag-6"),
		IndexName:              aws.String("TagIndex"),
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
		TableName: aws.String("article-tag-6"),
		Key: map[string]types.AttributeValue{
			"PK#PUB": &types.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s", item.Username, item.Publication)},
			"SK":     &types.AttributeValueMemberS{Value: item.Tag},
		},
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

	existingTags, err := t.Get(ctx, item)
	if err != nil {
		return nil, err
	}

	queryInput := dynamodb.QueryInput{
		TableName:              aws.String("article-tag-6"),
		KeyConditionExpression: aws.String("#v1 = :v1"),
		// FilterExpression:       aws.String("#v2 > :v2 AND NOT (Tag IN :v3)"),
		FilterExpression: aws.String("#v2 > :v2"),
		ExpressionAttributeNames: map[string]string{
			"#v1": "PK#PUB",
			"#v2": "TotalCount",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v1": &types.AttributeValueMemberS{Value: fmt.Sprintf("PUB#%s", item.Publication)},
			":v2": &types.AttributeValueMemberN{Value: strconv.Itoa(constant.PopularTagCount)},
		},
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
	filterExpression := fmt.Sprintf("%s AND NOT (Tag IN (", *queryInput.FilterExpression)

	filterAttr := []string{}

	for k, val := range existingTags {

		//
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

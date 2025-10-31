package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	"github.com/danilobml/user-manager/internal/user/model"
)

type UserRepositoryDdb struct {
	client    *dynamodb.Client
	tableName string
}

func NewUserRepositoryDdb(ddbClient *dynamodb.Client) *UserRepositoryDdb {
	return &UserRepositoryDdb{
		client:    ddbClient,
		tableName: "users",
	}
}

func (ur *UserRepositoryDdb) List(ctx context.Context) ([]*model.User, error) {
	out, err := ur.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(ur.tableName),
	})
	if err != nil {
		return nil, err
	}

	var users []model.User
	err = attributevalue.UnmarshalListOfMaps(out.Items, &users)
	if err != nil {
		return nil, err
	}

	var usersResp []*model.User
	for _, user := range users {
		usersResp = append(usersResp, &user)
	}

	return usersResp, nil
}

func (ur *UserRepositoryDdb) FindById(ctx context.Context, id uuid.UUID) (*model.User, error) {
	out, err := ur.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(ur.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, nil
	}

	var user model.User
	err = attributevalue.UnmarshalMap(out.Item, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepositoryDdb) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	out, err := ur.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(ur.tableName),
		IndexName:              aws.String("email-index"),
		KeyConditionExpression: aws.String("#email = :email"),
		ExpressionAttributeNames: map[string]string{
			"#email": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}
	if len(out.Items) == 0 {
		return nil, nil
	}
	var user model.User
	if err := attributevalue.UnmarshalMap(out.Items[0], &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepositoryDdb) Create(ctx context.Context, user model.User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = ur.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:                aws.String(ur.tableName),
		Item:                     item,
		ConditionExpression:      aws.String("attribute_not_exists(#id)"),
		ExpressionAttributeNames: map[string]string{"#id": "id"},
	})
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepositoryDdb) Update(ctx context.Context, user model.User) error {

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	_, err = ur.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(ur.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: user.ID.String()},
		},
		UpdateExpression:    aws.String("SET #email=:email, #hashed_password=:hashed_password, #roles=:roles"),
		ConditionExpression: aws.String("attribute_exists(#id)"),
		ExpressionAttributeNames: map[string]string{
			"#id":              "id",
			"#email":           "email",
			"#hashed_password": "hashed_password",
			"#roles":           "roles",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email":           item["email"],
			":hashed_password": item["hashed_password"],
			":roles":           item["roles"],
		},
	})

	return err
}

func (ur *UserRepositoryDdb) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := ur.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(ur.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

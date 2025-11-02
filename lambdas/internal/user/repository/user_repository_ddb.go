package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	dtos "github.com/danilobml/user-manager/internal/user/dtos"
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

	var ddbUsers []dtos.UserDDB
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &ddbUsers); err != nil {
		return nil, err
	}

	usersResp := make([]*model.User, 0, len(ddbUsers))
	for i := range ddbUsers {
		u, convErr := dtos.FromDDB(ddbUsers[i])
		if convErr != nil {
			return nil, convErr
		}
		user := u
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

	var ddbUser dtos.UserDDB
	if err := attributevalue.UnmarshalMap(out.Item, &ddbUser); err != nil {
		return nil, err
	}
	u, convErr := dtos.FromDDB(ddbUser)
	if convErr != nil {
		return nil, convErr
	}

	return &u, nil
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
	var ddbUser dtos.UserDDB
	if err := attributevalue.UnmarshalMap(out.Items[0], &ddbUser); err != nil {
		return nil, err
	}
	u, convErr := dtos.FromDDB(ddbUser)
	if convErr != nil {
		return nil, convErr
	}
	return &u, nil
}

func (ur *UserRepositoryDdb) Create(ctx context.Context, user model.User) error {
	ddbUser := dtos.ToDDB(user)
	item, err := attributevalue.MarshalMap(ddbUser)
	if err != nil {
		return err
	}
	item["id"] = &types.AttributeValueMemberS{Value: user.ID.String()}

	_, err = ur.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(ur.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(#id)"),
		ExpressionAttributeNames: map[string]string{
			"#id": "id",
		},
	})
	return err
}

func (ur *UserRepositoryDdb) Update(ctx context.Context, user model.User) error {
	ddbUser := dtos.ToDDB(user)
	av, err := attributevalue.MarshalMap(ddbUser)
	if err != nil { return fmt.Errorf("failed to marshal user: %w", err) }

	names := map[string]string{
		"#id":              "id",
		"#hashed_password": "hashed_password",
		"#roles":           "roles",
		"#is_active":       "is_active",
	}
	values := map[string]types.AttributeValue{
		":hashed_password": av["hashed_password"],
		":roles":           av["roles"],
		":is_active":       av["is_active"],
	}
	setParts := []string{"#hashed_password=:hashed_password", "#roles=:roles", "#is_active=:is_active"}

	if ddbUser.Email != "" {
		names["#email"] = "email"
		values[":email"] = av["email"]
		setParts = append(setParts, "#email=:email")
	}

	updateExpr := "SET " + strings.Join(setParts, ", ")

	_, err = ur.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(ur.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: ddbUser.ID},
		},
		UpdateExpression:          aws.String(updateExpr),
		ConditionExpression:       aws.String("attribute_exists(#id)"),
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
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
	return err
}

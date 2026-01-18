package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"myproject/internal/models"
)

type ConfigStore interface {
	GetConfig(ctx context.Context, deviceID string) (*models.UserConfig, error)
	PutConfig(ctx context.Context, config *models.UserConfig) error
	QueryByAlarmTime(ctx context.Context, dateBucket string, startTime, endTime string) ([]models.UserConfig, error)
}

type DynamoDBStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBStore(client *dynamodb.Client) *DynamoDBStore {
	return &DynamoDBStore{
		client:    client,
		tableName: "UserConfigs",
	}
}

var ErrConfigNotFound = errors.New("config not found")

func (s *DynamoDBStore) GetConfig(ctx context.Context, deviceID string) (*models.UserConfig, error) {
	result, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"deviceId": &types.AttributeValueMemberS{Value: deviceID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	if result.Item == nil {
		return nil, ErrConfigNotFound
	}

	var config models.UserConfig
	err = attributevalue.UnmarshalMap(result.Item, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (s *DynamoDBStore) PutConfig(ctx context.Context, config *models.UserConfig) error {
	item, err := attributevalue.MarshalMap(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	log.Printf("PutConfig: device=%s, enabled=%v, nextAlarmTime=%q, alarmDateBucket=%q", config.DeviceID, config.Enabled, config.NextAlarmTime, config.AlarmDateBucket)
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("DynamoDB PutItem error: %v", err)
		return fmt.Errorf("failed to put config: %w", err)
	}

	return nil
}

func (s *DynamoDBStore) QueryByAlarmTime(ctx context.Context, dateBucket string, startTime, endTime string) ([]models.UserConfig, error) {
	keyCondition := "alarmDateBucket = :bucket AND nextAlarmTime BETWEEN :start AND :end"

	input := &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		IndexName:              aws.String("AlarmTimeIndex"),
		KeyConditionExpression: aws.String(keyCondition),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":bucket": &types.AttributeValueMemberS{Value: dateBucket},
			":start":  &types.AttributeValueMemberS{Value: startTime},
			":end":    &types.AttributeValueMemberS{Value: endTime},
		},
	}

	result, err := s.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to query alarms: %w", err)
	}

	var configs []models.UserConfig
	err = attributevalue.UnmarshalListOfMaps(result.Items, &configs)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configs: %w", err)
	}

	return configs, nil
}

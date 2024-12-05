package repository

import (
	"bonded/internal/models"
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func (r *eventRepository) CreateEvent(ctx context.Context, calendar *models.Calendar, event *models.Event) error {
	calendar.Events = append(calendar.Events, *event)

	item, err := dynamodbattribute.MarshalMap(event)
	if err != nil {
		return err
	}

	item["CalendarID"] = &dynamodb.AttributeValue{S: aws.String(calendar.CalendarID)}
	item["SortKey"] = &dynamodb.AttributeValue{S: aws.String("EVENT#" + event.EventID)}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	}

	_, err = r.dynamoDB.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}

	gsiItem := map[string]*dynamodb.AttributeValue{
		"CalendarID": {S: aws.String(calendar.CalendarID)},
		"SortKey":    {S: aws.String("CAL#" + calendar.CalendarID + "#" + event.EventID)},
		"UserID":     {S: aws.String(calendar.OwnerUserID)},
	}
	gsiInput := &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      gsiItem,
	}

	_, err = r.dynamoDB.PutItemWithContext(ctx, gsiInput)
	return err
}

func (r *eventRepository) FindEvents(ctx context.Context, calendarID string) ([]*models.Event, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("CalendarID = :calendarID AND begins_with(SortKey, :sortPrefix)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":calendarID": {S: aws.String(calendarID)},
			":sortPrefix": {S: aws.String("EVENT#")},
		},
	}

	result, err := r.dynamoDB.QueryWithContext(ctx, input)
	if err != nil {
		return nil, err
	}

	events := make([]*models.Event, 0, len(result.Items))
	for _, item := range result.Items {
		var event models.Event
		err = dynamodbattribute.UnmarshalMap(item, &event)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}

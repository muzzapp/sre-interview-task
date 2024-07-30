package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type WriteRecordsInput struct {
	Records []WriteRecordsInputRecord
}

type WriteRecordsInputRecord struct {
	Component   string
	Environment string
	State       string
	Timestamp   time.Time
}

// Writes status records
func (r *Repo) WriteRecords(ctx context.Context, input *WriteRecordsInput) error {
	for _, record := range input.Records {
		// Write the current record
		ddbCurrentRecord := dynamoCurrentRecord{
			Component:   record.Component,
			Environment: record.Environment,
			Ts:          record.Timestamp,
			State:       record.State,
			Type:        RecordTypeCurrent,
		}
		avc, err := attributevalue.MarshalMap(ddbCurrentRecord)
		if err != nil {
			return fmt.Errorf("marshal map: %w", err)
		}
		_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(r.tableName),
			Item:      avc,
		})
		if err != nil {
			return fmt.Errorf("put item: %w", err)
		}

		// Write the historical record
		ddbHistoricalRecord := dynamoHistoricalRecord{
			// Re-use the timestamp as a record ID
			RecordId:    record.Timestamp.String(),
			Component:   record.Component,
			Environment: record.Environment,
			Ts:          record.Timestamp,
			State:       record.State,
			// This should use `RecordTypeCurrent` instead of hardcoding `"historical"`
			Type: "historical",
		}
		avh, err := attributevalue.MarshalMap(ddbHistoricalRecord)
		if err != nil {
			return fmt.Errorf("marshal map: %w", err)
		}
		_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(r.tableName),
			Item:      avh,
		})
		if err != nil {
			return fmt.Errorf("put item: %w", err)
		}
	}
	return nil
}

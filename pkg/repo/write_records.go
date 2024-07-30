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
			// All information is lost when this error is returned, logging the failed record
			// or returning the failed records to the client would help troubleshooting problems.
			return fmt.Errorf("marshal map: %w", err)
		}
		// Using PutItem results in a lot of calls compared to
		// using a TransactWriteItems call.
		// Because this is being done in multiple PutItem calls
		// an error may be returned to the client in a half
		// written state. This may not be a problem but should
		// be documented.
		// If a record with an earlier timestamp is delayed
		// it will overwrite a later record as there is no
		// constraint on the timestamp.
		_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(r.tableName),
			Item:      avc,
		})
		if err != nil {
			// Multiple `"put item"` and `"marshal map"` errors in this function may lead to confusion
			// on where the error originates from.
			return fmt.Errorf("put item: %w", err)
		}

		// Write the historical record
		ddbHistoricalRecord := dynamoHistoricalRecord{
			// Re-use the timestamp as a record ID
			// If two timestamps are the same for a component data will be
			// overwritten.
			RecordId:    record.Timestamp.String(),
			Component:   record.Component,
			Environment: record.Environment,
			Ts:          record.Timestamp,
			State:       record.State,
			// Use constant instead of hardcoded value
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

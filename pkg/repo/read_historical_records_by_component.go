package repo

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ReadHistoricalRecordsByComponentInput struct {
	Component string
}

type ReadHistoricalRecordsByComponentOutput struct {
	Records []ReadHistoricalRecordsByComponentOutputRecord
}

type ReadHistoricalRecordsByComponentOutputRecord struct {
	Component   string
	Environment string
	Timestamp   time.Time
	State       string
}

// Returns a list of "historical" records for a specific component
// This is used to analyse the historical of a deployment for auditing
// and incident investigation processes.
func (r *Repo) ReadHistoricalRecordsByComponent(_ context.Context, input *ReadHistoricalRecordsByComponentInput) (ReadHistoricalRecordsByComponentOutput, error) {
	res, err := r.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return ReadHistoricalRecordsByComponentOutput{}, fmt.Errorf("scan table: %w", err)
	}
	var records []ReadHistoricalRecordsByComponentOutputRecord
	for _, item := range res.Items {
		var dhr dynamoHistoricalRecord
		attributevalue.UnmarshalMap(item, &dhr)
		if dhr.Type == RecordTypeHistorical && dhr.Component == input.Component {
			records = append(records, ReadHistoricalRecordsByComponentOutputRecord{
				Component:   dhr.Component,
				Environment: dhr.Environment,
				Timestamp:   dhr.Ts,
				State:       dhr.State,
			})
		}
	}
	// Sort the records by time
	slices.SortFunc(records, func(i, j ReadHistoricalRecordsByComponentOutputRecord) int {
		if i.Timestamp.Before(j.Timestamp) {
			return -1
		}
		if i.Timestamp.After(j.Timestamp) {
			return 1
		}
		return 0
	})
	return ReadHistoricalRecordsByComponentOutput{Records: records}, nil
}

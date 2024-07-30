package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type GetCurrentStatusInput struct{}

type GetCurrentStatusOutput struct {
	Records []GetCurrentStatusOutputRecord
}

type GetCurrentStatusOutputRecord struct {
	Component   string
	Environment string
	Timestamp   time.Time
	State       string
}

// Returns a list of "current" records for all components/environments.
// This gives an overview of the current state of the world to identify any
// currently ongoing deployments or issues with current deployments.
func (r *Repo) GetCurrentStatus(_ context.Context, _ *GetCurrentStatusInput) (GetCurrentStatusOutput, error) {
	res, err := r.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return GetCurrentStatusOutput{}, fmt.Errorf("scan table: %w", err)
	}
	var records []GetCurrentStatusOutputRecord
	for _, item := range res.Items {
		var dcr dynamoCurrentRecord
		attributevalue.UnmarshalMap(item, &dcr)
		if dcr.Type == "current" {
			records = append(records, GetCurrentStatusOutputRecord{
				Component:   dcr.Component,
				Environment: dcr.Environment,
				Timestamp:   dcr.Ts,
				State:       dcr.State,
			})
		}
	}
	return GetCurrentStatusOutput{Records: records}, nil
}

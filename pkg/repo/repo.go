package repo

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	// All of these errors are actually unused
	ErrUnknownResource       = errors.New("unknown resource")
	ErrInvalidResourceFormat = errors.New("invalid resource format")
	ErrInvalidDataFormat     = errors.New("invalid data format")
)

type Repo struct {
	client    *dynamodb.Client
	tableName string
}

func New(client *dynamodb.Client, tableName string) *Repo {
	return &Repo{
		client:    client,
		tableName: tableName,
	}
}

const (
	RecordTypeCurrent    = "current"
	RecordTypeHistorical = "historical"
)

// The current state of the world based on the last record
type dynamoCurrentRecord struct {
	// Using long names for DynamoDB attributes consumes more RCU and WCUs
	// so they should be kept as short as possible.
	// See the table restructure comment in main.tf
	Component   string
	Environment string
	Ts          time.Time `dynamodbav:",unixtime"`
	State       string
	Type        string
}

// The historical state of the world
type dynamoHistoricalRecord struct {
	// We need the record ID to not be a primary key
	// so we-use the Component attribute and rename
	// Component to Cmp to avoid conflicts.
	RecordId    string `dynamodbav:"Component,unixtime"`
	Environment string
	Component   string    `dynamodbav:"Cmp"`
	Ts          time.Time `dynamodbav:",unixtime"`
	State       string
	Type        string
}

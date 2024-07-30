package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	auditrepo "github.com/muzzapp/interviewtask/pkg/repo"
)

func main() {
	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("load aws config", "error", err)
		return
	}
	ddbClient := dynamodb.NewFromConfig(awsConfig)
	repo := auditrepo.New(ddbClient, "cicd-audit")

	// A deployment started
	err = repo.WriteRecords(ctx, &auditrepo.WriteRecordsInput{
		Records: []auditrepo.WriteRecordsInputRecord{
			{
				Component:   "monitoring.pipeline.router",
				Environment: "dev",
				State:       "deploy-started",
				Timestamp:   time.Now().Add(-time.Second * 10),
			},
			{
				Component:   "monitoring.pipeline.router",
				Environment: "dev",
				State:       "deploy-failed",
				Timestamp:   time.Now(),
			},
			{
				Component:   "social.api",
				Environment: "dev-rev0",
				State:       "deploy-started",
				Timestamp:   time.Now(),
			},
			{
				Component:   "social.api",
				Environment: "dev-rev0",
				State:       "deploy-success",
				Timestamp:   time.Now().Add(time.Second * 5),
			},
		},
	})
	if err != nil {
		slog.Error("write records", "error", err)
		return
	}
	fmt.Println("")
	fmt.Println("Current status overview:")
	fmt.Println("")
	resp, err := repo.GetCurrentStatus(ctx, &auditrepo.GetCurrentStatusInput{})
	if err != nil {
		slog.Error("read current", "error", err)
		return
	}
	for _, record := range resp.Records {
		fmt.Printf("%v %s %s %s\n", record.Timestamp, record.Component, record.Environment, record.State)
	}
	fmt.Println("")
	fmt.Println("Historical status for monitoring.pipeline.router:")
	fmt.Println("")
	historicalResp, err := repo.ReadHistoricalRecordsByComponent(ctx, &auditrepo.ReadHistoricalRecordsByComponentInput{Component: "monitoring.pipeline.router"})
	if err != nil {
		slog.Error("read current", "error", err)
		return
	}
	for _, record := range historicalResp.Records {
		fmt.Printf("%v %s %s %s\n", record.Timestamp, record.Component, record.Environment, record.State)
	}
}

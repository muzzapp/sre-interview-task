package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/muzzapp/interviewtask/pkg/repo"
)

func main() {
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		slog.Error("load aws config", "error", err)
		return
	}
	ddbClient := dynamodb.NewFromConfig(awsConfig)
	repo := repo.New(ddbClient, "cicd-audit")
	h := handler{auditrepo: repo}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /current", h.getCurrentStatus)
	mux.HandleFunc("GET /historical/{component}", h.getHistoricalStatusByComponent)
	mux.HandleFunc("POST /records", h.postNewRecords)
	slog.Info("starting api", "address", ":8080")
	http.ListenAndServe(":8080", mux)
}

type handler struct {
	auditrepo *repo.Repo
}

func (h *handler) getCurrentStatus(w http.ResponseWriter, _ *http.Request) {
	resp, err := h.auditrepo.GetCurrentStatus(context.Background(), &repo.GetCurrentStatusInput{})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	dat, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(dat)
}

func (h *handler) getHistoricalStatusByComponent(w http.ResponseWriter, r *http.Request) {
	resp, err := h.auditrepo.ReadHistoricalRecordsByComponent(context.Background(), &repo.ReadHistoricalRecordsByComponentInput{
		Component: r.PathValue("component"),
	})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	dat, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(dat)
}

func (h *handler) postNewRecords(w http.ResponseWriter, r *http.Request) {
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	var input repo.WriteRecordsInput
	err = json.Unmarshal(dat, &input)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err = h.auditrepo.WriteRecords(context.Background(), &input)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Ok"))
}

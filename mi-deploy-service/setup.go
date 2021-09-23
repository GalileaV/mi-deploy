package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/option"
)

var (
	triggersService *cloudbuild.ProjectsTriggersService
	kgKey           string
	slackSecret     string
	projectId       string
	triggerId       string
)

func setup(ctx context.Context) {
	kgKey = os.Getenv("GCP_API_KEY")
	slackSecret = os.Getenv("SLACK_SIGN_IN_SECRET")
	projectId = os.Getenv("PROJECT_ID")
	triggerId = os.Getenv("TRIGGER_ID")

	if triggersService == nil {
		cloudbuildService, err := cloudbuild.NewService(ctx, option.WithAPIKey(kgKey))
		if err != nil {
			log.Fatalf("cloudbuild.NewService: %v", err)
		}
		triggersService = cloudbuild.NewProjectsTriggersService(cloudbuildService)
	}
}

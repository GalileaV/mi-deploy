package mideploy

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/option"
)

var (
	triggersService *cloudbuild.ProjectsTriggersService
	slackSecret     string
	projectId       string
	triggerId       string
	environment     string
	scopes          []string
	aud             string
)

func setup(ctx context.Context) {
	slackSecret = os.Getenv("SLACK_SIGN_IN_SECRET")
	projectId = os.Getenv("PROJECT_ID")
	triggerId = os.Getenv("TRIGGER_ID")
	environment = os.Getenv("ENVIRONMENT")

	prvKey, err := ioutil.ReadFile("serverless_function_source_code/misalud-development-cb.json")
	if err != nil {
		log.Fatalln(err)
	}

	scopes = []string{"https://www.googleapis.com/auth/cloud-platform"}
	aud = "https://oauth2.googleapis.com/token"

	conf, err := google.JWTConfigFromJSON(prvKey, scopes...)
	if err != nil {
		log.Fatalf("Unable to generate token source: %v", err)
	}

	conf.Audience = aud

	if triggersService == nil {
		cloudbuildService, err := cloudbuild.NewService(ctx, option.WithTokenSource(conf.TokenSource(ctx)))
		if err != nil {
			log.Fatalf("cloudbuild.NewService: %v", err)
		}
		triggersService = cloudbuild.NewProjectsTriggersService(cloudbuildService)
	}
}

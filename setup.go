package mideploy

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var (
	cloudbuildClient *cloudbuild.Client
	slackSecret      string
	projectId        string
	triggerId        string
	environment      string
	scopes           []string
	aud              string
)

func setup(ctx context.Context) {
	slackSecret = os.Getenv("SLACK_SIGN_IN_SECRET")
	projectId = os.Getenv("PROJECT_ID")
	triggerId = os.Getenv("TRIGGER_ID")
	environment = os.Getenv("ENVIRONMENT")

	path := fmt.Sprintf("serverless_function_source_code/misalud-%s-cb.json", environment)
	prvKey, err := ioutil.ReadFile(path)
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

	if cloudbuildClient == nil {
		cloudbuildClient, err = cloudbuild.NewClient(ctx, option.WithTokenSource(conf.TokenSource(ctx)))
		if err != nil {
			log.Fatalf("cloudbuild.NewService: %v", err)
		}
	}
}

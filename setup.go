package mideploy

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/option"
)

var (
	triggersService *cloudbuild.ProjectsTriggersService
	slackSecret     string
	projectId       string
	triggerId       string
)

func setup(r *http.Request) {
	slackSecret = os.Getenv("SLACK_SIGN_IN_SECRET")
	projectId = os.Getenv("PROJECT_ID")
	triggerId = os.Getenv("TRIGGER_ID")

	googleOauthConfig := &oauth2.Config{
		ClientID:     os.Getenv("0AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("0AUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/cloud-platform"},
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://oauth2.googleapis.com/token",
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		},
	}

	ctx := r.Context()
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	code := r.FormValue("code")
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	if triggersService == nil {
		cloudbuildService, err := cloudbuild.NewService(ctx, option.WithTokenSource(googleOauthConfig.TokenSource(ctx, token)))
		if err != nil {
			log.Fatalf("cloudbuild.NewService: %v", err)
		}
		triggersService = cloudbuild.NewProjectsTriggersService(cloudbuildService)
	}
}

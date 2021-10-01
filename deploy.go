package mideploy

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	cloudbuildpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

const (
	version                     = "v0"
	slackRequestTimestampHeader = "X-Slack-Request-Timestamp"
	slackSignatureHeader        = "X-Slack-Signature"
)

type attachment struct {
	Color     string `json:"color"`
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	ImageURL  string `json:"image_url"`
}

// Message is the a Slack message event.
// see https://api.slack.com/docs/message-formatting
type Message struct {
	ResponseType string       `json:"response_type"`
	Text         string       `json:"text"`
	Attachments  []attachment `json:"attachments"`
}

func MiDeploy(w http.ResponseWriter, r *http.Request) {
	setup(r.Context())

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Couldn't read request body: %v", err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Couldn't parse form", 400)
		log.Fatalf("ParseForm: %v", err)
	}

	// Reset r.Body as ParseForm depletes it by reading the io.ReadCloser.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if r.Form["team_domain"][0] == "misaludworkspace" {
		slackSecret = os.Getenv("SLACK_SIGN_IN_SECRET")
	} else {
		slackSecret = os.Getenv("EXTERNAL_SIGN_IN_SECRET")
	}

	result, err := verifyWebHook(r, slackSecret)
	if err != nil {
		log.Fatalf("verifyWebhook: %v", err)
	}
	if !result {
		log.Fatalf("signatures did not match.")
	}

	if len(r.Form["text"]) == 0 {
		log.Fatalf("empty text in form")
	}

	response_url = r.Form["response_url"][0]
	branchName = r.Form["text"][0]

	message := &Message{
		ResponseType: "ephemeral",
		Text:         "_Request received_ :cat-typing:",
	}
	sendSlackMessage(message)
	runTrigger(strings.Title(r.Form["user_name"][0]))
}

func sendSlackMessage(message *Message) {
	client := &http.Client{Timeout: 10 * time.Second}
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error Occurred jsonData. %+v", err)
	}

	req, err := http.NewRequest(http.MethodPost, response_url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error Occurred with slack request, line 104. %+v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to Slack endpoint. %+v", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}
	log.Println("Slack message, response Body:", string(body))
}

func runTrigger(author string) {
	t := time.Now()
	defer func() {
		t2 := time.Since(t)
		log.Println("the runTrigger function take", t2)
	}()

	ctx := context.Background()
	values := cloudbuildpb.RepoSource_BranchName{
		BranchName: branchName,
	}
	reposource := cloudbuildpb.RepoSource{
		Revision: &values,
	}
	req := &cloudbuildpb.RunBuildTriggerRequest{
		ProjectId: projectId,
		TriggerId: triggerId,
		Source:    &reposource,
	}

	resp, err := cloudbuildClient.RunBuildTrigger(ctx, req)
	if err != nil {
		message := &Message{
			ResponseType: "ephemeral",
			Text:         fmt.Sprintf(":warning: Something went wrong trying to run the trigger. :warning:\nWe couldn't find the *%s* branch. Please, review the spelling and try again.", branchName),
		}
		sendSlackMessage(message)
		log.Printf("Error running trigger: %v", err)
	} else if resp == nil {
		log.Printf("empty response")
	} else {
		message, err := formatSlackMessage(author)
		if err != nil {
			log.Println("message unavailable")
		}
		sendSlackMessage(message)
		log.Printf("Trigger response %v", resp)
	}
}

// verifyWebHook verifies the request signature.
// See https://api.slack.com/docs/verifying-requests-from-slack.
func verifyWebHook(r *http.Request, slackSigningSecret string) (bool, error) {
	timeStamp := r.Header.Get(slackRequestTimestampHeader)
	slackSignature := r.Header.Get(slackSignatureHeader)

	if timeStamp == "" || slackSignature == "" {
		return false, fmt.Errorf("either timeStamp or signature headers were blank")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return false, fmt.Errorf("ioutil.ReadAll(%v): %v", r.Body, err)
	}

	// Reset the body so other calls won't fail.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	baseString := fmt.Sprintf("%s:%s:%s", version, timeStamp, body)

	signature := getSignature([]byte(baseString), []byte(slackSigningSecret))

	trimmed := strings.TrimPrefix(slackSignature, fmt.Sprintf("%s=", version))
	signatureInHeader, err := hex.DecodeString(trimmed)

	if err != nil {
		return false, fmt.Errorf("hex.DecodeString(%v): %v", trimmed, err)
	}

	return hmac.Equal(signature, signatureInHeader), nil
}

func getSignature(base []byte, secret []byte) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write(base)

	return h.Sum(nil)
}

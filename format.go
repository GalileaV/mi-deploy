package mideploy

import (
	"fmt"

	"google.golang.org/api/cloudbuild/v1"
)

func formatSlackMessage(branchName string, author string, response *cloudbuild.Operation) (*Message, error) {
	if response == nil {
		return nil, fmt.Errorf("empty response")
	}
	attach := attachment{Color: "#3367d6", Text: fmt.Sprintf("%s is deploying %s to %s :deploy_now:", author, branchName, environment)}

	message := &Message{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("branchName: %s", branchName),
		Attachments:  []attachment{attach},
	}
	return message, nil
}

// [END functions_slack_format]

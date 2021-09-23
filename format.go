package mideploy

import (
	"fmt"

	"google.golang.org/api/cloudbuild/v1"
)

func formatSlackMessage(branchName string, response *cloudbuild.Operation) (*Message, error) {
	if response == nil {
		return nil, fmt.Errorf("empty response")
	}

	if !response.Done {
		message := &Message{
			ResponseType: "in_channel",
			Text:         fmt.Sprintf("branchName: %s", branchName),
			Attachments: []attachment{
				{
					Color: "#d6334b",
					Text:  "Deployment in progress.",
				},
			},
		}
		return message, nil
	}
	attach := attachment{Color: "#3367d6", Text: string(response.Response)}

	message := &Message{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("branchName: %s", branchName),
		Attachments:  []attachment{attach},
	}
	return message, nil
}

// [END functions_slack_format]

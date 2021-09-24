package mideploy

import (
	"fmt"
	"strings"

	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
)

func formatSlackMessage(branchName string, author string, response *longrunningpb.Operation) (*Message, error) {
	if response == nil {
		return nil, fmt.Errorf("empty response")
	}
	environmentTitle := strings.ToUpper(environment)
	attach := attachment{Color: "#3367d6", Text: fmt.Sprintf("*%s* is deploying %s to %s :deploy_now:", author, branchName, environmentTitle)}

	message := &Message{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("%s is using *%s*", author, environmentTitle),
		Attachments:  []attachment{attach},
	}
	return message, nil
}

package model

import "time"

const DefaultTimeout = 10 * time.Second

type SuccessResponse struct {
	State       string `json:"state"`
	Description string `json:"description"`
	Context     string `json:"context"`
	TargetURL   string `json:"target_url"`
}

type slackMessageBlock struct {
	Section string        `json:"type"`
	Text    slackMarkdown `json:"text"`
}

type slackMarkdown struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	Verbatim bool   `json:"verbatim"`
}

type slackMessage struct {
	Blocks []slackMessageBlock `json:"blocks"`
}

func CreateSlackMessage(message string) *slackMessage {
	return &slackMessage{
		Blocks: []slackMessageBlock{{
			Section: "section",
			Text: slackMarkdown{
				Type:     "mrkdwn",
				Text:     message,
				Verbatim: true,
			}},
		},
	}
}

// ErrMessage is a message that is returned on error.
type ErrMessage struct {
	Errors []string `json:"errors"`
}

// JSONError Function returns creates vector of errors from multiple errors.
func JSONError(errors ...error) ErrMessage {
	ret := ErrMessage{Errors: make([]string, len(errors))}
	for i, err := range errors {
		ret.Errors[i] = err.Error()
	}
	return ret
}

type pushData struct {
	PushedAt float64 `json:"pushed_at"`
	Tag      string  `json:"tag"`
	Pusher   string  `json:"pusher"`
}

type repoData struct {
	Name string `json:"repo_name"`
}

type DockerhubPayload struct {
	CallbackURL string   `json:"callback_url"`
	PushData    pushData `json:"push_data"`
	Repository  repoData `json:"repository"`
}

type harborRepository struct {
	Repository   string `json:"name"`
	NameSpace    string `json:"namespace"`
	FullRepoName string `json:"repo_full_name"`
}

type harborResource struct {
	Tag      string `json:"tag"`
	Resource string `json:"resource_url"`
}

type harborEvent struct {
	Resources  []harborResource `json:"resources"`
	Repository harborRepository `json:"repository"`
}

type HarborPayload struct {
	EventType string      `json:"type"`
	Time      float64     `json:"occur_at"`
	User      string      `json:"operator"`
	EventData harborEvent `json:"event_data"`
}

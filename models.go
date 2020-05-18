package main

type successResponse struct {
	State       string `json:"state"`
	Description string `json:"description"`
	Context     string `json:"context"`
	TargetURL   string `json:"target_url"`
}

type slackRequestBody struct {
	Text string `json:"text"`
}

// ErrMessage is a message that is returned on error
type ErrMessage struct {
	Errors []string `json:"errors"`
}

// JSONError Function returns creates vector of errors from multiple errors
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

type dockerhubPayload struct {
	CallbackURL string   `json:"callback_url"`
	PushData    pushData `json:"push_data"`
	Repository  repoData `json:"repository"`
}

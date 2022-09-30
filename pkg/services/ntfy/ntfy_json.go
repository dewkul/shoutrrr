package ntfy

import "fmt"

// messageRequest is the actual payload being sent to the Ntfy API
type messageRequest struct {
	Topic    string   `json:"topic"`
	Message  string   `json:"message,omitempty"`
	Title    string   `json:"title,omitempty"`
	Tags     []string `json:"tag,omitempty"`
	Priority uint8    `json:"priority,omitempty"`
	// TODO: Action Buttons
	Click    string `json:"click,omitempty"`
	Attach   string `json:"attach,omitempty"`
	FileName string `json:"filename,omitempty"`
	Delay    string `json:"delay,omitempty"`
	Email    string `json:"email,omitempty"`
}

type messageResponse struct {
	messageRequest
	ID        string `json:"id"`
	Timestamp uint64 `json:"timestamp"`
	Event     string `json:"event"`
}

type errorResponse struct {
	Name        string `json:"error"`
	Code        uint64 `json:"errorCode"`
	Description string `json:"errorDescription"`
}

func (er *errorResponse) Error() string {
	return fmt.Sprintf("server responds with %v (%v): %v", er.Name, er.Code, er.Description)
}

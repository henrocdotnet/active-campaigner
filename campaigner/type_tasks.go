package campaigner

// RequestContactTagCreate holds a JSON compatible request for creating contact tags.
// This is what is sent to the API for creation.
type RequestContactTagCreate struct {
	ContactID int64 `json:"contact"`
	TagID int64 `json:"tag"`
}
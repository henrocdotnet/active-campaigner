package campaigner

type TaskContactTagCreate struct {
	ContactID int64 `json:"contact"`
	TagID int64 `json:"tag"`
}
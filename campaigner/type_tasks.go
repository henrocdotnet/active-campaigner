package campaigner

// RequestContactUpdate holds a JSOn compatible request for updating contacts.
type RequestContactUpdate struct {
	ID             int64  `json:"id,omitempty"`
	EmailAddress   string `json:"email,omitempty"`
	FirstName      string `json:"firstName,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	PhoneNumber    string `json:"phone,omitempty"`
	IsDeleted      bool   `json:"deleted,omitempty"`
	OrganizationID int64  `json:"orgid,omitempty"`
}

// RequestContactTagCreate holds a JSON compatible request for creating contact tags.
// This is what is sent to the API for creation.
type RequestContactTagCreate struct {
	ContactID int64 `json:"contact"`
	TagID     int64 `json:"tag"`
}

// RequestContactFieldUpdate holds a JSON compatible request for updating contact custom fields.
type RequestContactFieldUpdate struct {
	ContactID int64  `json:"contact"`
	FieldID   int64  `json:"field"`
	Value     string `json:"value"`
}

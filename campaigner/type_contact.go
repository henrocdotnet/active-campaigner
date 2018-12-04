package campaigner

// TODO(organization): Should probably move these back into contact.go.

// This was generated from JSON returned by a read call.
type Contact struct {
	ID           int64  `json:"id"`
	EmailAddress string `json:"email"`
	PhoneNumber  string `json:"phone"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	OrgID        string `json:"orgid"`

	DateCreated string `json:"cdate"`
	DateUpdated string `json:"udate"`
	DateDeleted string `json:"deleted_at"`
	IsDeleted   int    `json:"deleted,string"`

	// These fields are still in progress.
	SegmentioID         string        `json:"segmentio_id"`
	BouncedHard         string        `json:"bounced_hard"`
	BouncedSoft         string        `json:"bounced_soft"`
	BouncedDate         interface{}   `json:"bounced_date"`
	IP                  string        `json:"ip"`
	Ua                  interface{}   `json:"ua"`
	Hash                string        `json:"hash"`
	SocialdataLastcheck interface{}   `json:"socialdata_lastcheck"`
	EmailLocal          string        `json:"email_local"`
	EmailDomain         string        `json:"email_domain"`
	Sentcnt             string        `json:"sentcnt"`
	RatingTstamp        interface{}   `json:"rating_tstamp"`
	Gravatar            string        `json:"gravatar"`
	Anonymized          string        `json:"anonymized"`
	Adate               interface{}   `json:"adate"`
	Edate               interface{}   `json:"edate"`
	CreatedUtcTimestamp string        `json:"created_utc_timestamp"`
	UpdatedUtcTimestamp string        `json:"updated_utc_timestamp"`
	ContactAutomations  []interface{} `json:"contactAutomations"`
	ContactLists        []interface{} `json:"contactLists"`
	FieldValues         []interface{} `json:"fieldValues"`
	GeoIps              []interface{} `json:"geoIps"`
	Deals               []interface{} `json:"deals"`
	Links               ContactLinks  `json:"links"`
	Organization        interface{}   `json:"organization"`
}

// Not sure what these link to at this point (other than the obvious).
type ContactLinks struct {
	BounceLogs         string `json:"bounceLogs"`
	ContactAutomations string `json:"contactAutomations"`
	ContactData        string `json:"contactData"`
	ContactGoals       string `json:"contactGoals"`
	ContactLists       string `json:"contactLists"`
	ContactLogs        string `json:"contactLogs"`
	ContactTags        string `json:"contactTags"`
	ContactDeals       string `json:"contactDeals"`
	Deals              string `json:"deals"`
	FieldValues        string `json:"fieldValues"`
	GeoIps             string `json:"geoIps"`
	Notes              string `json:"notes"`
	Organization       string `json:"organization"`
	PlusAppend         string `json:"plusAppend"`
	TrackingLogs       string `json:"trackingLogs"`
	ScoreValues        string `json:"scoreValues"`
}

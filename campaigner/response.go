package campaigner

// ResponseContactCreate holds a JSON compatible response for creating contacts.
type ResponseContactCreate struct {
	Contact Contact `json:"contact"`
}

// ResponseContactList holds a JSON compatible response for listing contacts.
type ResponseContactList struct {
	ScoreValues []interface{} `json:"scoreValues"`
	Contacts    []Contact     `json:"contacts"`
	Meta        struct {
		Total     string `json:"total"`
		PageInput struct {
			Segmentid  int         `json:"segmentid"`
			Formid     int         `json:"formid"`
			Listid     int         `json:"listid"`
			Tagid      int         `json:"tagid"`
			Limit      int         `json:"limit"`
			Offset     int         `json:"offset"`
			Search     interface{} `json:"search"`
			Sort       interface{} `json:"sort"`
			Seriesid   int         `json:"seriesid"`
			Waitid     int         `json:"waitid"`
			Status     int         `json:"status"`
			ForceQuery int         `json:"forceQuery"`
			Cacheid    string      `json:"cacheid"`
		} `json:"page_input"`
	} `json:"meta"`
	/*
	Contacts    []struct {
		Cdate               string        `json:"cdate"`
		Email               string        `json:"email"`
		Phone               string        `json:"phone"`
		FirstName           string        `json:"firstName"`
		LastName            string        `json:"lastName"`
		OrgID               string        `json:"orgid"`
		SegmentioID         string        `json:"segmentio_id"`
		BouncedHard         string        `json:"bounced_hard"`
		BouncedSoft         string        `json:"bounced_soft"`
		BouncedDate         string        `json:"bounced_date"`
		IP                  string        `json:"ip"`
		Ua                  string        `json:"ua"`
		Hash                string        `json:"hash"`
		SocialDataLastCheck string        `json:"socialdata_lastcheck"`
		EmailLocal          string        `json:"email_local"`
		EmailDomain         string        `json:"email_domain"`
		Sentcnt             string        `json:"sentcnt"`
		RatingTstamp        string        `json:"rating_tstamp"`
		Gravatar            string        `json:"gravatar"`
		Deleted             string        `json:"deleted"`
		Anonymized          string        `json:"anonymized"`
		Adate               string        `json:"adate"`
		Udate               string        `json:"udate"`
		Edate               string        `json:"edate"`
		DeletedAt           string        `json:"deleted_at"`
		CreatedUtcTimestamp string        `json:"created_utc_timestamp"`
		UpdatedUtcTimestamp string        `json:"updated_utc_timestamp"`
		ScoreValues         []interface{} `json:"scoreValues"`
		Links               ContactLinks  `json:"links"`
		ID                  string        `json:"id"`
		Organization        string        `json:"organization"`
	} `json:"contacts"`
	*/
}

// ResponseContactRead holds a JSON compatible response for reading contacts.
type ResponseContactRead struct {
	Contact            Contact `json:"contact"`
	ContactAutomations []struct {
		Contact           string      `json:"contact"`
		Seriesid          string      `json:"seriesid"`
		Startid           string      `json:"startid"`
		Status            string      `json:"status"`
		Adddate           string      `json:"adddate"`
		Remdate           interface{} `json:"remdate"`
		Timespan          interface{} `json:"timespan"`
		Lastblock         string      `json:"lastblock"`
		Lastdate          string      `json:"lastdate"`
		CompletedElements string      `json:"completedElements"`
		TotalElements     string      `json:"totalElements"`
		Completed         int         `json:"completed"`
		CompleteValue     int         `json:"completeValue"`
		Links             struct {
			Automation   string `json:"automation"`
			Contact      string `json:"contact"`
			ContactGoals string `json:"contactGoals"`
		} `json:"links"`
		ID         string `json:"id"`
		Automation string `json:"automation"`
	} `json:"contactAutomations"`
	ContactLists []struct {
		Contact               string      `json:"contact"`
		List                  string      `json:"list"`
		Form                  interface{} `json:"form"`
		Seriesid              string      `json:"seriesid"`
		Sdate                 interface{} `json:"sdate"`
		Udate                 interface{} `json:"udate"`
		Status                string      `json:"status"`
		Responder             string      `json:"responder"`
		Sync                  string      `json:"sync"`
		Unsubreason           interface{} `json:"unsubreason"`
		Campaign              interface{} `json:"campaign"`
		Message               interface{} `json:"message"`
		FirstName             string      `json:"first_name"`
		LastName              string      `json:"last_name"`
		IP4Sub                string      `json:"ip4Sub"`
		Sourceid              string      `json:"sourceid"`
		AutosyncLog           interface{} `json:"autosyncLog"`
		IP4Last               string      `json:"ip4_last"`
		IP4Unsub              string      `json:"ip4Unsub"`
		UnsubscribeAutomation interface{} `json:"unsubscribeAutomation"`
		Links                 struct {
			Automation            string `json:"automation"`
			List                  string `json:"list"`
			Contact               string `json:"contact"`
			Form                  string `json:"form"`
			AutosyncLog           string `json:"autosyncLog"`
			Campaign              string `json:"campaign"`
			UnsubscribeAutomation string `json:"unsubscribeAutomation"`
			Message               string `json:"message"`
		} `json:"links"`
		ID         string      `json:"id"`
		Automation interface{} `json:"automation"`
	} `json:"contactLists"`
	Deals []struct {
		Owner        string      `json:"owner"`
		Contact      string      `json:"contact"`
		Organization interface{} `json:"organization"`
		Group        interface{} `json:"group"`
		Title        string      `json:"title"`
		Nexttaskid   string      `json:"nexttaskid"`
		Currency     string      `json:"currency"`
		Status       string      `json:"status"`
		Links        struct {
			Activities   string `json:"activities"`
			Contact      string `json:"contact"`
			ContactDeals string `json:"contactDeals"`
			Group        string `json:"group"`
			NextTask     string `json:"nextTask"`
			Notes        string `json:"notes"`
			Organization string `json:"organization"`
			Owner        string `json:"owner"`
			ScoreValues  string `json:"scoreValues"`
			Stage        string `json:"stage"`
			Tasks        string `json:"tasks"`
		} `json:"links"`
		ID       string      `json:"id"`
		NextTask interface{} `json:"nextTask"`
	} `json:"deals"`
	FieldValues []struct {
		Contact string      `json:"contact"`
		Field   string      `json:"field"`
		Value   interface{} `json:"value"`
		Cdate   string      `json:"cdate"`
		Udate   string      `json:"udate"`
		Links   struct {
			Owner string `json:"owner"`
			Field string `json:"field"`
		} `json:"links"`
		ID    string `json:"id"`
		Owner string `json:"owner"`
	} `json:"fieldValues"`
	// TODO(JSON): This field isn't being sent at all at the moment (null).
	GeoAddresses []struct {
		IP4      string        `json:"ip4"`
		Country2 string        `json:"country2"`
		Country  string        `json:"country"`
		State    string        `json:"state"`
		City     string        `json:"city"`
		Zip      string        `json:"zip"`
		Area     string        `json:"area"`
		Lat      string        `json:"lat"`
		Lon      string        `json:"lon"`
		Tz       string        `json:"tz"`
		Tstamp   string        `json:"tstamp"`
		Links    []interface{} `json:"links"`
		ID       string        `json:"id"`
	} `json:"geoAddresses"`
	GeoIps []struct {
		Contact    string `json:"contact"`
		Campaignid string `json:"campaignid"`
		Messageid  string `json:"messageid"`
		Geoaddrid  string `json:"geoaddrid"`
		IP4        string `json:"ip4"`
		Tstamp     string `json:"tstamp"`
		GeoAddress string `json:"geoAddress"`
		Links      struct {
			GeoAddress string `json:"geoAddress"`
		} `json:"links"`
		ID string `json:"id"`
	} `json:"geoIps"`
}

// ResponseError holds a list of ActiveCampaign errors.
type ResponseError struct {
	// TODO: Not in use, what was I doing here?
	Errors []ActiveCampaignError
}

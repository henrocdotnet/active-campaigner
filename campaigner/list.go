package campaigner

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ListContactAdd adds a contact to a list.
func (c *Campaigner) ListContactAdd(listID int64, contactID int64) (response ResponseListContactAdd, err error) {
	// Check that both the contact and list exist.
	_, err = c.ContactRead(contactID)
	if err != nil {
		return response, fmt.Errorf("list contact addition failed, could not find contact: %s", err)
	}
	l, err := c.ListRead(listID)
	if err != nil {
		return response, fmt.Errorf("list contact addition failed, could not find list: %s", err)
	}

	req := RequestListContactAdd{ListID: listID, ContactID: contactID, Status: true}

	u := "/api/3/contactLists"
	r, body, err := c.post(u, map[string]interface{}{"contactList": req})

	// Response check.
	switch r.StatusCode {
	case http.StatusOK, http.StatusCreated:
		if err = json.Unmarshal(body, &response); err != nil {
			return response, fmt.Errorf("list contact addition failed, JSON error: %s (%s)", err, string(body))
		}

		response.Custom.ListName = l.List.Name

		return response, nil
	default:
		return response, fmt.Errorf("list contact addition failed, unspecified error (%d): %s", r.StatusCode, string(body))
	}
}

// ListList lists available contact lists.
func (c *Campaigner) ListList() (response ResponseListList, err error) {
	// Send GET request.
	u := "/api/3/lists"
	r, body, err := c.get(u)

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		if err = json.Unmarshal(body, &response); err != nil {
			return response, fmt.Errorf("list listing failed, JSON error: %s", err)
		}

		return response, nil
	default:
		return response, fmt.Errorf("list listing failed, unspecified error (%d): %s", r.StatusCode, string(body))
	}
}

// ListRead reads a contact list.
func (c *Campaigner) ListRead(id int64) (response ResponseListRead, err error) {
	// Send GET request.
	u := fmt.Sprintf("/api/3/lists/%d", id)
	r, body, err := c.get(u)

	// Response check.
	switch r.StatusCode {
	case http.StatusOK:
		if err = json.Unmarshal(body, &response); err != nil {
			return response, fmt.Errorf("list read failed, JSON error: %s", err)
		}

		return response, nil
	default:
		return response, fmt.Errorf("list read failed, unspecified error (%d): %s", r.StatusCode, string(body))
	}
}

// List holds a JSON compatible list as it exists in the API.
type List struct {
	ID                   int64       `json:"id,string"`
	Name                 string      `json:"name"`
	AnalyticsDomains     interface{} `json:"analytics_domains"`
	AnalyticsSource      string      `json:"analytics_source"`
	AnalyticsUa          string      `json:"analytics_ua"`
	CarbonCopy           interface{} `json:"carboncopy"`
	DateCreated          string      `json:"cdate"`
	DateDeleted          interface{} `json:"deletestamp"`
	FacebookSession      interface{} `json:"facebook_session"`
	FullAddress          string      `json:"fulladdress"`
	GetUnsubscribeReason string      `json:"get_unsubscribe_reason"`
	Links                struct {
		AddressLists     string `json:"addressLists"`
		ContactGoalLists string `json:"contactGoalLists"`
		User             string `json:"user"`
	} `json:"links"`
	OptInMessageID       string      `json:"optinmessageid"`
	OptInOptOut          string      `json:"optinoptout"`
	OptOutConfig         string      `json:"optoutconf"`
	PEmbedImage          string      `json:"p_embed_image"`
	PUseAnalyticsLink    string      `json:"p_use_analytics_link"`
	PUseAnalyticsRead    string      `json:"p_use_analytics_read"`
	PUseCaptcha          string      `json:"p_use_captcha"`
	PUseFacebook         string      `json:"p_use_facebook"`
	PUseTracking         string      `json:"p_use_tracking"`
	PUseTwitter          string      `json:"p_use_twitter"`
	Private              string      `json:"private"`
	RequireName          string      `json:"require_name"`
	SendLastBroadcast    string      `json:"send_last_broadcast"`
	SenderAddr1          string      `json:"sender_addr1"`
	SenderAddr2          string      `json:"sender_addr2"`
	SenderCity           string      `json:"sender_city"`
	SenderCountry        string      `json:"sender_country"`
	SenderName           string      `json:"sender_name"`
	SenderPhone          string      `json:"sender_phone"`
	SenderReminder       string      `json:"sender_reminder"`
	SenderState          string      `json:"sender_state"`
	SenderURL            string      `json:"sender_url"`
	SenderZip            string      `json:"sender_zip"`
	StringID             string      `json:"stringid"`
	SubscriptionNotify   interface{} `json:"subscription_notify"`
	ToName               string      `json:"to_name"`
	TwitterToken         string      `json:"twitter_token"`
	TwitterTokenSecret   string      `json:"twitter_token_secret"`
	DateUpdated          interface{} `json:"udate"`
	UnsubscriptionNotify interface{} `json:"unsubscription_notify"`
	User                 string      `json:"user"`
	UserID               string      `json:"userid"`
}

// RequestListContactAdd holds a JSON compatible request for adding contacts to lists.
type RequestListContactAdd struct {
	ListID    int64 `json:"list"`
	ContactID int64 `json:"contact"`
	Status    bool  `json:"status"`
}

// ResponseListContactAdd holds a JSON compatible response for adding contacts to lists.
type ResponseListContactAdd struct {
	Custom struct {
		ListName string
	}
	ContactList struct {
		Automation  interface{} `json:"automation"`
		AutoSyncLog interface{} `json:"autosyncLog"`
		Campaign    interface{} `json:"campaign"`
		ContactID   Int64json   `json:"contact"`
		FirstName   string      `json:"first_name"`
		Form        interface{} `json:"form"`
		ID          string      `json:"id"`
		IP4Sub      string      `json:"ip4Sub"`
		IP4Unsub    string      `json:"ip4Unsub"`
		IP4Last     string      `json:"ip4_last"`
		LastName    string      `json:"last_name"`
		Links       struct {
			Automation            string `json:"automation"`
			AutoSyncLog           string `json:"autosyncLog"`
			Campaign              string `json:"campaign"`
			Contact               string `json:"contact"`
			Form                  string `json:"form"`
			List                  string `json:"list"`
			Message               string `json:"message"`
			UnsubscribeAutomation string `json:"unsubscribeAutomation"`
		} `json:"links"`
		ListID                Int64json   `json:"list"`
		Message               interface{} `json:"message"`
		Responder             string      `json:"responder"`
		DateSubscribed        string      `json:"sdate"`
		SeriesID              Int64json   `json:"seriesid"`
		SourceID              Int64json   `json:"sourceid"`
		Status                int         `json:"status"`
		Sync                  string      `json:"sync"`
		UnsubscribeReason     string      `json:"unsubreason"`
		UnsubscribeAutomation interface{} `json:"unsubscribeAutomation"`
	} `json:"contactList"`
	Contacts []Contact `json:"contacts"`
}

// ResponseListList holds a json compatible response for listing contact lists.
type ResponseListList struct {
	Lists []List `json:"lists"`
	Meta  struct {
		Total string `json:"total"`
	} `json:"meta"`
}

// ResponseListRead holds a JSON compatible struct for reading lists.
type ResponseListRead struct {
	List List `json:"list"`
}

package campaigner

import (
	"fmt"
	"log"
	"strings"
)

// ActiveCampaignError holds a JSON compatible API error (nested structure, see ResponseError).
type ActiveCampaignError struct {
	// TODO(move): me.
	Errors []ActiveCampaignErrorLine
}

// Error satisfies the error interface.  Generates and returns the error string.
func (e ActiveCampaignError) Error() string {
	var o []string

	for x, y := range e.Errors {
		o = append(o, fmt.Sprintf(`<error index="%d" title="%s" detail="%s" code="%s" />`, x, y.Title, y.Detail, y.Code))
	}

	return fmt.Sprintf("ActiveCampaign error: %s", strings.Join(o, ", "))
}

// ActiveCampaignErrorLine holds a JSON compatible API error line (nested structure, see ActiveCampaignError).
type ActiveCampaignErrorLine struct {
	Code   string                    `json:"code"`
	Detail string                    `json:"detail"`
	Source ActiveCampaignErrorSource `json:"source"`
	Title  string                    `json:"title"`
}

// ActiveCampaignErrorSource holds a JSON compatible API error source (nested structure, see ActiveCampaignErrorLine).
type ActiveCampaignErrorSource struct {
	Pointer string `json:"pointer"`
}

// ActiveCampaignErrorList holds a JSON compatible list of API error lines (nested structure, see ActiveCampaignErrorLine).
type ActiveCampaignErrorList struct {
	// TODO(api): Remove this (probably, unused at the moment, try to think back).
	List []ActiveCampaignErrorLine `json:"errors"`
}

// CustomError holds a group of HTTP errors during an API call.  Allows more than one error to be returned at a time.
type CustomError struct {
		// TODO(naming): Rename this.
	HTTPErrors []error
	Message    string
}

// Error satisfies the error interface.  Generates and returns the error string.
func (e CustomError) Error() string {
	var l []string

	for _, y := range e.HTTPErrors {
		l = append(l, y.Error())
	}
	return fmt.Sprintf("%s (%s)", e.Message, strings.Join(l, ", "))
}

// SetMessage sets the error message string (sprintf style).
func (e CustomError) SetMessage(m string, a ...interface{}) CustomError {
	e.Message = fmt.Sprintf(m, a...)
	return e
}

// WriteToLog logs the contents of an error.
func (e CustomError) WriteToLog() {
	var (
		list   []string
		output string
	)

	for _, x := range e.HTTPErrors {
		list = append(list, x.Error())
	}

	output = fmt.Sprintf("campaigner error: %s:\n%s", e.Message, strings.Join(list, "\n"))

	log.Printf(output)
}

// CustomErrorNotFound is an error subtype that allows for a specific condition to be checked for.
type CustomErrorNotFound struct {
	CustomError
}

package campaigner

import (
	"fmt"
	"log"
	"strings"
)

type ActiveCampaignError struct {
	Errors []ActiveCampaignErrorLine
}

func (e ActiveCampaignError) Error() string {
	var o []string

	for x, y := range e.Errors {
		o = append(o, fmt.Sprintf(`<error index="%d" title="%s" detail="%s" code="%s" />`, x, y.Title, y.Detail, y.Code))
	}

	return fmt.Sprintf("ActiveCampaign error: %s", strings.Join(o, ", "))
}

type ActiveCampaignErrorLine struct {
	Code   string                    `json:"code"`
	Detail string                    `json:"detail"`
	Source ActiveCampaignErrorSource `json:"source"`
	Title  string                    `json:"title"`
}

type ActiveCampaignErrorSource struct {
	Pointer string `json:"pointer"`
}

type ActiveCampaignErrorList struct {
	List []ActiveCampaignErrorLine `json:"errors"`
}

// TODO(naming): Rename this.
type CustomError struct {
	HTTPErrors []error
	Message    string
}

func (e CustomError) Error() string {
	var l []string

	for _, y := range e.HTTPErrors {
		l = append(l, y.Error())
	}
	return fmt.Sprintf("%s (%s)", e.Message, strings.Join(l, ", "))
}

func (e CustomError) SetMessage(m string) CustomError {
	e.Message = m
	return e
}

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

type CustomErrorNotFound struct {
	CustomError
}

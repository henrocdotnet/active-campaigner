# Active Campaigner

Provides interaction with [ActiveCampaign](https://www.activecampaign.com/) API.

This project is still in a very early rough draft phase.

# Basic Usage 
Calls are named _Type_ + _Action_ so that they are grouped together.  I've tried to name methods so that they are
self-explanatory.

## Create Contact
```go
c := campaigner.Campaigner{ ApiToken: "token", BaseURL: "url" }
contact := campaigner.Contact{ FirstName: "First", LastName: "Last", EmailAddress: "first.last@domain.com" }
response, _ := c.ContactCreate(contact) // error handling omitted for brevity
log.Printf("API response data (type ResponseContactCreate) in response variable: %#v\n", response)
```

## Create Organization
```go
c := campaigner.Campaigner{ ApiToken: "token", BaseURL: "url" }
org := campaigner.Organization{ Name: "Org" }
response, _ := c.OrganizationCreate(org) // error handling omitted for brevity
log.Printf("API response data (type ResponseOrganizationCreate) in response variable: %#v\n", response)
```

# Unit Test Setup

## Config
```bash
export AC_API_TOKEN='your token goes here'
export AC_BASE_URL='https://your-subdomain.api-us1.com'
export AC_UNIT_TEST_EMAIL='address@domain.com'
export AC_UNIT_TEST_PHONE='2125551212'
```

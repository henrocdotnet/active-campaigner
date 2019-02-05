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

# API Bugs
* Contact Update: Using PUT method will not update an existing contact.  I have not tried using PUT to create a new contact yet. ([forum link](https://community.activecampaign.com/t/possible-bug-v3-contact-update-put-attempts-failed-with-email-exists/5961))
* Contact Update: The Organization ID returned in the contact JSON is sometimes a string and sometimes an int.  This appears
to depend whether the organization ID is sent in the request JSON.

Current User
    "contact":{id: "35", "email":"test@user.com","firstName":"Test","lastName":"User","phone":"2125551212"}}

Update Variants
    {"contact":{"email":"test@user.com","firstName":"Test","lastName":"Smith","phone":"7185551212"}} // Last and phone update, get email exists.
    {"contact":{"email":"test@user.com","firstName":"Test","lastName":"User","phone":"2125551212","orgid":1}} // Add OrgID only, get email exists.
    {"contact":{"id":35,"email":"test@user.com","firstName":"Test","lastName":"User","phone":"2125551212","orgid":1}} // Include ID, get email exists.
    {"contact":{"firstName":"Test","lastName":"Smith","phone":"7185551212"}} // URL includes ID, get email required.
    {"contact":{"id":35,"firstName":"Test","lastName":"Smith","phone":"7185551212","orgid":1}} // JSON also includes ID, get email required.


# Active Campaigner

Provides interaction with [ActiveCampaign](https://www.activecampaign.com/) API.

This project is still in a very early rough draft phase.

# Basic Usage
## Create Contact
```
c := campaigner.Campaigner{ ApiToken: "token", BaseURL: "url" }
contact := campaigner.Contact{ FirstName: "First", LastName: "Last", EmailAddress: "first.last@domain.com" }
response, err := campaigner.ContactCreate(contact)
if err != nil {
    // Handle error.
}
log.Printf("API response data (type ResponseContactCreate) in response variable: %#v\n", response)

```

## Create Organization
```
c := campaigner.Campaigner{ ApiToken: "token", BaseURL: "url" }
org := campaigner.Organization{ Name: "Org" }
response, err := campaigner.OrganizationCreate(org)
if err != nil {
    // Handle error.
}
log.Printf("API response data (type ResponseOrganizationCreate) in response variable: %#v\n", response)

```

<img src="https://henroc.net/semita/" />

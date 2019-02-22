package main


const classTemplate = `
<?php


class ActiveCampaignerTagConstants
{

{{- range $x, $y := .Tags }}

    // {{ $y.Name }}
    const {{ $y.Name | cleanTagName }} = {{ $y.ID }};

{{- end }}

}
`

const fieldsTemplate = `
<?php


class ActiveCampaignerFieldConstants
{

{{- range $x, $y := .Fields }}

    // {{ $y.Title }}
    const {{ $y.Title | cleanTagName }} = {{ $y.ID }};

{{- end }}

}
`

const listsTemplate = `
<?php


class ActiveCampaignerListConstants
{

{{- range $x, $y := .Lists }}

    // {{ $y.Name }}
    const {{ $y.Name | cleanTagName }} = {{ $y.ID }};

{{- end }}

}
`

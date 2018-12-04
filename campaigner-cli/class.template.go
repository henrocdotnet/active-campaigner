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

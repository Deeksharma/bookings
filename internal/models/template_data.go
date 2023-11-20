package models

import "github.com/Deeksharma/bookings/internal/forms"

// it will contain all data base models template data models

// ypu inevitably want some data to be rendered in every template, so you'll always want to add the some data to template data

// TemplateData holds data sent from handlers to templates - holds information for every single page in our wensite
type TemplateData struct {
	StringData map[string]string
	IntData    map[string]int32
	FloatData  map[string]float32
	Data       map[string]interface{}
	CSRFToken  string // cross site request forgery token, used in post forms
	// these data will be sent back to the client
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
}

// js mein {{.CSRFToken}} will git value
// in go templates use {{index .StringData "key_string"}}

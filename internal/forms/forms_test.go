package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) { // if the function has receiver then this is the naming convention
	postedData := url.Values{}
	form := NewForm(postedData)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) { // if the function has receiver then this is the naming convention
	postedData := url.Values{}
	form := NewForm(postedData)

	form.Required("a", "b", "c")

	isValid := form.Valid()
	if isValid {
		t.Error("got valid when should have been invalid")
	}

	postedData = url.Values{}
	postedData.Add("a", "-")
	postedData.Add("b", "-")
	postedData.Add("c", "-")

	form = NewForm(postedData)
	form.Required("a", "b", "c")
	isValid = form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := NewForm(postedData)

	has := form.Has("d")
	if has {
		t.Error("form shows value is present even if it doesnt")
	}

	postedData = url.Values{}
	postedData.Add("a", "cmek")
	postedData.Add("b", "-")
	postedData.Add("c", "-")

	form = NewForm(postedData)
	has = form.Has("a")
	if !has {
		t.Error("shows form does not has field when it should have")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := NewForm(postedData)

	form.IsEmail("email-2")
	isValid := form.Valid()
	if isValid {
		t.Error("shows valid when email non existent")
	}

	postedData = url.Values{}
	postedData.Add("email-2", "a@a.com")
	form = NewForm(postedData)

	form.IsEmail("email-2")
	isValid = form.Valid()
	if !isValid {
		t.Error("shows invalid when email is valid")
	}

	postedData = url.Values{}
	postedData.Add("email-1", "fdge")
	form = NewForm(postedData)
	form.IsEmail("email-1")
	isValid = form.Valid()
	if isValid {
		t.Error("shows valid when email is not valid")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := NewForm(postedData)

	form.MinLength("x", 10)
	isValid := form.Valid()
	if isValid {
		t.Error("form shows min length for non existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error but didn't get one")
	}

	postedData = url.Values{}
	postedData.Add("a", "abc")

	form = NewForm(postedData)

	form.MinLength("a", 11)
	isValid = form.Valid()
	if isValid {
		t.Error("shows minlength 11 when lenght is 3")
	}

	postedData = url.Values{}
	postedData.Add("a", "abc")

	form = NewForm(postedData)

	form.MinLength("a", 3)
	isValid = form.Valid()
	if !isValid {
		t.Error("shows minlength 11 when lenght is 3")
	}
	isError = form.Errors.Get("a")
	if isError != "" {
		t.Error("should not have an error but shows get one")
	}
}

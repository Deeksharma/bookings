package forms

type errors map[string][]string

// Add adds an  error message for a particular type
func (e errors) Add(field, message string){
	e[field] = append(e[field], message)
}

// Get return the first error message on the field
func (e errors) Get(field string) string{
	if len(e[field]) == 0 {
		return ""
	}
	return e[field][0]
}
package output

import "errors"

type Template struct {
	template string
}

// Create string from template with input parameters
func (template *Template) Process(field string, metric string, payload map[string]string) (error, string) {
	return errors.New("not yet implemented"), ""
}

func NewTempate(template string) (error, *Template) {
	var err error

	t := new(Template)

	if err = validateTemplate(template); err != nil {
		return err, nil
	}

	return nil, t
}

func validateTemplate(s string) error {

	return nil
}


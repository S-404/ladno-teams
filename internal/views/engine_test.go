package views

import "testing"

func TestModalTemplatesParsed(t *testing.T) {
	e, err := NewEngine()
	if err != nil {
		t.Fatal(err)
	}
	names := []string{
		"partials/admin/modals/invite",
		"partials/admin/modals/team",
		"partials/admin/modals/workspace",
		"partials/admin/modals/workspace_role",
	}
	for _, name := range names {
		if e.templates.Lookup(name) == nil {
			t.Errorf("template %q not parsed", name)
		}
	}
}

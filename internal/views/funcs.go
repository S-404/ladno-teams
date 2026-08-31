package views

import (
	"html/template"
	"ladno-teams/internal/entity"
	"time"
)

const adminDateFormat = "2006-01-02 15:04:05"

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"formatTime":    formatTime,
		"strVal":        strVal,
		"userLabel":     userLabel,
		"checkbox":      checkboxData,
		"teamSelect":    teamSelectData,
		"sectionHeader": sectionHeader,
		"modalFormHTMX": modalFormHTMX,
	}
}

func formatTime(t time.Time) string {
	return t.Format(adminDateFormat)
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func userLabel(u entity.AdminUserView) string {
	if u.Name != nil && *u.Name != "" {
		return *u.Name + " (" + u.Login + ")"
	}
	return u.Login
}

func checkboxData(postURL, fieldName string, checked bool) CheckboxForm {
	value := "false"
	if checked {
		value = "true"
	}
	return CheckboxForm{
		PostURL:   postURL,
		FieldName: fieldName,
		Value:     value,
		Checked:   checked,
	}
}

func teamSelectData(teams []entity.Team, selected string) TeamSelect {
	return TeamSelect{Teams: teams, Selected: selected}
}

func sectionHeader(title, actionLabel, actionURL string) SectionHeader {
	return SectionHeader{
		Title:       title,
		ActionLabel: actionLabel,
		ActionURL:   actionURL,
	}
}

func modalFormHTMX() template.HTMLAttr {
	return template.HTMLAttr(AdminModalFormHTMX)
}

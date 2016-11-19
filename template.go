package cronv

import (
	"fmt"
	"text/template"
	"time"
)

const TEMPLATE = "vis.tpl"

func MakeTemplate() *template.Template {
	funcMap := template.FuncMap{
		"CronvIter": func(cronv *Cronv) <-chan *Exec {
			return cronv.Iter()
		},
		"JSEscapeString": func(v string) string {
			return template.JSEscapeString(v)
		},
		"NewJsDate": func(v time.Time) string {
			return fmt.Sprintf("new Date(%d,%d,%d,%d,%d)", v.Year(), v.Month(), v.Day(), v.Hour(), v.Minute())
		},
		"DateFormat": func(v time.Time, format string) string {
			return v.Format(format)
		},
	}
	return template.Must(template.New(TEMPLATE).Funcs(funcMap).ParseFiles(TEMPLATE))
}

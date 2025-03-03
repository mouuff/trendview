package provider

import "time"

type Report struct {
	Title    string
	Content  string
	DateTime *time.Time
	Link     string
	GUID     string
}

type ReportProvider interface {
	GetReports() ([]Report, error)
}

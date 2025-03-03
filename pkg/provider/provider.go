package provider

type Report struct {
	Title    string
	Content  string
	DateTime string
	Link     string
	GUID     string
}

type ReportProvider interface {
	GetReports() ([]Report, error)
}

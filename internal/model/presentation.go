package model

type Job struct {
	Role, Company string
}
type ContactCommercial struct {
	Name, Email, Phone string
}
type Experience struct {
	Firm                string
	Role                string
	Date                string
	Description         string
	Tasks, Technologies []string
}
type CommercialProposition struct {
	Role           string
	ConsultantName string
	Localisation   string
	Availability   string
}
type Presentation struct {
	Job                   Job
	ExperienceSummary     string
	Experiences           []Experience
	CommercialProposition CommercialProposition
	ContactCommercial     ContactCommercial
}

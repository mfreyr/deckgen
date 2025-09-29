package models

// Job represents the job details from the job advertising.
type Job struct {
	Role, Company string
}

// ContactCommercial holds the contact information for the commercial representative.
type ContactCommercial struct {
	Name, Email, Phone string
}

// Experience details a single professional experience.
type Experience struct {
	Firm          string
	Role          string
	Date          string
	Description   string
	Tasks         []string
	Technologies  []string
}

// CommercialProposition contains the commercial aspects of the proposal.
type CommercialProposition struct {
	Role           string
	ConsultantName string
	Localisation   string
	Availability   string
}

// Presentation is the root model, aggregating all information for the final output.
type Presentation struct {
	Job                   Job
	ExperienceSummary     string
	Experiences           []Experience
	CommercialProposition CommercialProposition
	ContactCommercial     ContactCommercial
}
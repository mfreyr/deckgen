package model

type File struct {
	ID        int
	Name      string
	Extension string
	Content   []byte
}

type JobAd struct {
	ID                      int      `json:"id"`
	Title                   string   `json:"title"`
	CompanyName             string   `json:"company_name"`
	Location                string   `json:"location"`
	KeyResponsibilities     []string `json:"key_responsibilities"`
	RequiredQualifications  []string `json:"required_qualifications"`
	PreferredQualifications []string `json:"preferred_qualifications"`
	RawText                 string   `json:"raw_text"`
}

type Experience struct {
	CompanyName string `json:"company_name"`
	Dates       string `json:"dates"`
	JobTitle    string `json:"job_title"`
	Description string `json:"description"`
	Tools       string `json:"tools"`
}

type CandidateResume struct {
	ID               int          `json:"id"`
	FullName         string       `json:"full_name"`
	Description      string       `json:"description"`
	ShortDescription string       `json:"short_description"`
	Experiences      []Experience `json:"experiences"`
	Certifications   []string     `json:"certifications"`
	Skills           []string     `json:"skills"`
	Location         string       `json:"location"`
	Availability     string       `json:"availability"`
	Facturation      string       `json:"facturation"`
	AverageDailyRate string       `json:"average_daily_rate"`
	BillingMode      string       `json:"billing_mode"`
}

type CandidateAdaptedResume struct {
	ID     int             `json:"id"`
	JobAd  JobAd           `json:"job_ad"`
	Resume CandidateResume `json:"resume"`
}

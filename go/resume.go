package main

type Resume struct {
	Basics       Basics        `yaml:"basics"`
	Disposition  Disposition   `yaml:"disposition"`
	Work         []WorkEntry   `yaml:"work"`
	Projects     []Project     `yaml:"projects"`
	Skills       Skills        `yaml:"skills"`
	Certificates []Certificate `yaml:"certificates"`
	Education    []Education   `yaml:"education"`
	Languages    []Language    `yaml:"languages"`
	Interests    []Interest    `yaml:"interests"`
	Testimonials []Testimonial `yaml:"testimonials"`
	References   []Reference   `yaml:"references"`
}

type Basics struct {
	Name     string    `yaml:"name"`
	Label    string    `yaml:"label"`
	Email    string    `yaml:"email"`
	Phone    string    `yaml:"phone"`
	URL      string    `yaml:"url"`
	Summary  string    `yaml:"summary"`
	Location Location  `yaml:"location"`
	Profiles []Profile `yaml:"profiles"`
}

type Location struct {
	City        string `yaml:"city"`
	Region      string `yaml:"region"`
	CountryCode string `yaml:"countryCode"`
}

type Profile struct {
	Network  string `yaml:"network"`
	Username string `yaml:"username"`
	URL      string `yaml:"url"`
}

type Disposition struct {
	Travel        int        `yaml:"travel"`
	Authorization string     `yaml:"authorization"`
	Commitment    []string   `yaml:"commitment"`
	Remote        bool       `yaml:"remote"`
	Relocation    Relocation `yaml:"relocation"`
}

type Relocation struct {
	Willing      bool     `yaml:"willing"`
	Destinations []string `yaml:"destinations"`
}

type WorkEntry struct {
	Employer      string   `yaml:"employer"`
	EmployerGroup string   `yaml:"employerGroup"`
	Position      string   `yaml:"position"`
	URL           string   `yaml:"url"`
	StartDate     string   `yaml:"startDate"`
	EndDate       string   `yaml:"endDate"`
	Summary       string   `yaml:"summary"`
	Location      string   `yaml:"location"`
	Highlights    []string `yaml:"highlights"`
	Keywords      []string `yaml:"keywords"`
}

type Project struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	URL         string   `yaml:"url"`
	Type        string   `yaml:"type"`
	Highlights  []string `yaml:"highlights"`
	Keywords    []string `yaml:"keywords"`
	Roles       []string `yaml:"roles"`
	StartDate   string   `yaml:"startDate"`
	EndDate     string   `yaml:"endDate"`
}

type Skills struct {
	Sets []SkillSet  `yaml:"sets"`
	List []SkillItem `yaml:"list"`
}

type SkillSet struct {
	Name   string   `yaml:"name"`
	Skills []string `yaml:"skills"`
}

type SkillItem struct {
	Name    string `yaml:"name"`
	Level   string `yaml:"level"`
	Summary string `yaml:"summary"`
	Years   int    `yaml:"years"`
}

type Certificate struct {
	Name             string `yaml:"name"`
	Date             string `yaml:"date"`
	URL              string `yaml:"url"`
	Issuer           string `yaml:"issuer"`
	ID               string `yaml:"id"`
	VerificationCode string `yaml:"verificationCode"`
}

type Education struct {
	Institution string `yaml:"institution"`
	URL         string `yaml:"url"`
	Area        string `yaml:"area"`
	StudyType   string `yaml:"studyType"`
	StartDate   string `yaml:"startDate"`
	EndDate     string `yaml:"endDate"`
	Score       string `yaml:"score"`
	Location    string `yaml:"location"`
}

type Language struct {
	Language string `yaml:"language"`
	Fluency  string `yaml:"fluency"`
	Years    int    `yaml:"years"`
}

type Interest struct {
	Name    string `yaml:"name"`
	Summary string `yaml:"summary"`
}

type Testimonial struct {
	Name     string `yaml:"name"`
	Role     string `yaml:"role"`
	Category string `yaml:"category"`
	URL      string `yaml:"url"`
	Email    string `yaml:"email"`
	Quote    string `yaml:"quote"`
}

type Reference struct {
	Name     string `yaml:"name"`
	Role     string `yaml:"role"`
	Category string `yaml:"category"`
	Email    string `yaml:"email"`
	URL      string `yaml:"url"`
}

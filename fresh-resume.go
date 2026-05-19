package main

type FRESHResume struct {
	Schema               string                `json:"$schema"`
	Title                string                `json:"title"`
	Type                 string                `json:"type"`
	AdditionalProperties bool                  `json:"additionalProperties"`
	Properties           FRESHResumeProperties `json:"properties"`
}

type FRESHResumeProperties struct {
	Name            Name            `json:"name"`
	Meta            Meta            `json:"meta"`
	Info            Info            `json:"info"`
	Disposition     Disposition     `json:"disposition"`
	Contact         Contact         `json:"contact"`
	Location        Location        `json:"location"`
	Employment      Employment      `json:"employment"`
	Projects        Projects        `json:"projects"`
	Skills          Skills          `json:"skills"`
	Service         Affiliation     `json:"service"`
	Education       Education       `json:"education"`
	Social          Social          `json:"social"`
	Recognition     Reading         `json:"recognition"`
	Writing         Writing         `json:"writing"`
	Reading         Reading         `json:"reading"`
	Speaking        Speaking        `json:"speaking"`
	Governance      Governance      `json:"governance"`
	Languages       Languages       `json:"languages"`
	Samples         Samples         `json:"samples"`
	References      References      `json:"references"`
	Testimonials    Testimonials    `json:"testimonials"`
	Interests       Interests       `json:"interests"`
	Extracurricular Extracurricular `json:"extracurricular"`
	Affiliation     Affiliation     `json:"affiliation"`
}

type Affiliation struct {
	Type                 string                `json:"type"`
	AdditionalProperties bool                  `json:"additionalProperties"`
	Description          string                `json:"description"`
	Properties           AffiliationProperties `json:"properties"`
}

type AffiliationProperties struct {
	Summary PhoneClass    `json:"summary"`
	History PurpleHistory `json:"history"`
}

type PurpleHistory struct {
	Type            CommitmentType  `json:"type"`
	AdditionalItems bool            `json:"additionalItems"`
	Items           GovernanceItems `json:"items"`
}

type GovernanceItems struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           PurpleProperties `json:"properties"`
}

type PurpleProperties struct {
	Category     PhoneClass  `json:"category"`
	Organization Name        `json:"organization"`
	Role         *PhoneClass `json:"role,omitempty"`
	URL          *EmailClass `json:"url,omitempty"`
	Start        EmailClass  `json:"start"`
	End          EmailClass  `json:"end"`
	Summary      PhoneClass  `json:"summary"`
	Highlights   Commitment  `json:"highlights"`
	Keywords     Commitment  `json:"keywords"`
	Location     *PhoneClass `json:"location,omitempty"`
	Position     *PhoneClass `json:"position,omitempty"`
}

type PhoneClass struct {
	Type        SummaryType `json:"type"`
	Description string      `json:"description"`
}

type EmailClass struct {
	Type        SummaryType `json:"type"`
	Description string      `json:"description"`
	Format      Format      `json:"format"`
}

type Commitment struct {
	Type            CommitmentType `json:"type"`
	Description     *string        `json:"description,omitempty"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           PhoneClass     `json:"items"`
	Required        *bool          `json:"required,omitempty"`
}

type Name struct {
	Type        SummaryType `json:"type"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Format      *Format     `json:"format,omitempty"`
}

type Contact struct {
	Type                 string            `json:"type"`
	AdditionalProperties bool              `json:"additionalProperties"`
	Description          string            `json:"description"`
	Properties           ContactProperties `json:"properties"`
}

type ContactProperties struct {
	Email   EmailClass `json:"email"`
	Phone   PhoneClass `json:"phone"`
	Website EmailClass `json:"website"`
	Other   Other      `json:"other"`
}

type Other struct {
	Type            CommitmentType `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           OtherItems     `json:"items"`
}

type OtherItems struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           FluffyProperties `json:"properties"`
}

type FluffyProperties struct {
	Label    PhoneClass `json:"label"`
	Category PhoneClass `json:"category"`
	Value    PhoneClass `json:"value"`
}

type Disposition struct {
	Type                 string                `json:"type"`
	AdditionalProperties bool                  `json:"additionalProperties"`
	Description          string                `json:"description"`
	Properties           DispositionProperties `json:"properties"`
}

type DispositionProperties struct {
	Travel        PhoneClass `json:"travel"`
	Authorization PhoneClass `json:"authorization"`
	Commitment    Commitment `json:"commitment"`
	Remote        PhoneClass `json:"remote"`
	Relocation    Relocation `json:"relocation"`
}

type Relocation struct {
	Type                 string               `json:"type"`
	AdditionalProperties bool                 `json:"additionalProperties"`
	Properties           RelocationProperties `json:"properties"`
}

type RelocationProperties struct {
	Willing      Years      `json:"willing"`
	Destinations Commitment `json:"destinations"`
}

type Years struct {
	Type        []string `json:"type"`
	Description string   `json:"description"`
}

type Education struct {
	Type                 string              `json:"type"`
	AdditionalProperties bool                `json:"additionalProperties"`
	Description          string              `json:"description"`
	Properties           EducationProperties `json:"properties"`
}

type EducationProperties struct {
	Summary PurpleSummary `json:"summary"`
	Level   Name          `json:"level"`
	Degree  PhoneClass    `json:"degree"`
	History FluffyHistory `json:"history"`
}

type FluffyHistory struct {
	Type            CommitmentType `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           PurpleItems    `json:"items"`
}

type PurpleItems struct {
	Type                 string              `json:"type"`
	AdditionalProperties bool                `json:"additionalProperties"`
	Properties           TentacledProperties `json:"properties"`
}

type TentacledProperties struct {
	Title       PhoneClass `json:"title"`
	Institution Name       `json:"institution"`
	Area        PhoneClass `json:"area"`
	StudyType   PhoneClass `json:"studyType"`
	Start       EmailClass `json:"start"`
	End         EmailClass `json:"end"`
	Grade       PhoneClass `json:"grade"`
	Curriculum  Commitment `json:"curriculum"`
	URL         EmailClass `json:"url"`
	Summary     PhoneClass `json:"summary"`
	Keywords    Commitment `json:"keywords"`
	Highlights  Commitment `json:"highlights"`
	Location    PhoneClass `json:"location"`
}

type PurpleSummary struct {
	Type        SummaryType `json:"type"`
	Description string      `json:"description:"`
}

type Employment struct {
	Type                 string               `json:"type"`
	Description          string               `json:"description"`
	AdditionalProperties bool                 `json:"additionalProperties"`
	Properties           EmploymentProperties `json:"properties"`
}

type EmploymentProperties struct {
	Summary PurpleSummary    `json:"summary"`
	History TentacledHistory `json:"history"`
}

type TentacledHistory struct {
	Type            CommitmentType `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           FluffyItems    `json:"items"`
}

type FluffyItems struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           StickyProperties `json:"properties"`
}

type StickyProperties struct {
	Employer   Name       `json:"employer"`
	Position   PhoneClass `json:"position"`
	URL        EmailClass `json:"url"`
	Start      EmailClass `json:"start"`
	End        EmailClass `json:"end"`
	Summary    PhoneClass `json:"summary"`
	Highlights Commitment `json:"highlights"`
	Location   PhoneClass `json:"location"`
	Keywords   Commitment `json:"keywords"`
}

type Extracurricular struct {
	Type            CommitmentType       `json:"type"`
	Description     string               `json:"description"`
	AdditionalItems bool                 `json:"additionalItems"`
	Items           ExtracurricularItems `json:"items"`
}

type ExtracurricularItems struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           IndigoProperties `json:"properties"`
}

type IndigoProperties struct {
	Title    Name       `json:"title"`
	Activity Name       `json:"activity"`
	Location PhoneClass `json:"location"`
	Start    EmailClass `json:"start"`
	End      EmailClass `json:"end"`
}

type Governance struct {
	Type            CommitmentType  `json:"type"`
	AdditionalItems bool            `json:"additionalItems"`
	Description     string          `json:"description"`
	Items           GovernanceItems `json:"items"`
}

type Info struct {
	Type                 string         `json:"type"`
	AdditionalProperties bool           `json:"additionalProperties"`
	Description          string         `json:"description"`
	Properties           InfoProperties `json:"properties"`
}

type InfoProperties struct {
	Label PhoneClass `json:"label"`
	Class PhoneClass `json:"class"`
	Image PhoneClass `json:"image"`
	Brief PhoneClass `json:"brief"`
	Quote PhoneClass `json:"quote"`
}

type Interests struct {
	Type            CommitmentType `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Description     string         `json:"description"`
	Items           InterestsItems `json:"items"`
}

type InterestsItems struct {
	Type                 string             `json:"type"`
	AdditionalProperties bool               `json:"additionalProperties"`
	Properties           IndecentProperties `json:"properties"`
}

type IndecentProperties struct {
	Name     Name          `json:"name"`
	Summary  FluffySummary `json:"summary"`
	Keywords Commitment    `json:"keywords"`
}

type FluffySummary struct {
	Type SummaryType `json:"type"`
}

type Languages struct {
	Type            CommitmentType `json:"type"`
	Description     string         `json:"description"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           LanguagesItems `json:"items"`
}

type LanguagesItems struct {
	Type                 string              `json:"type"`
	AdditionalProperties bool                `json:"additionalProperties"`
	Properties           HilariousProperties `json:"properties"`
}

type HilariousProperties struct {
	Language Name       `json:"language"`
	Level    PhoneClass `json:"level"`
	Years    Years      `json:"years"`
}

type Location struct {
	Type                 string             `json:"type"`
	Description          string             `json:"description"`
	AdditionalProperties bool               `json:"additionalProperties"`
	Properties           LocationProperties `json:"properties"`
}

type LocationProperties struct {
	Address PhoneClass `json:"address"`
	Code    PhoneClass `json:"code"`
	City    PhoneClass `json:"city"`
	Country PhoneClass `json:"country"`
	Region  PhoneClass `json:"region"`
}

type Meta struct {
	Type                 string         `json:"type"`
	AdditionalProperties bool           `json:"additionalProperties"`
	Required             bool           `json:"required"`
	Description          string         `json:"description"`
	Properties           MetaProperties `json:"properties"`
}

type MetaProperties struct {
	Format  Name       `json:"format"`
	Version PhoneClass `json:"version"`
}

type Projects struct {
	Type            CommitmentType `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Description     string         `json:"description"`
	Items           ProjectsItems  `json:"items"`
}

type ProjectsItems struct {
	Type                 string              `json:"type"`
	AdditionalProperties bool                `json:"additionalProperties"`
	Properties           AmbitiousProperties `json:"properties"`
}

type AmbitiousProperties struct {
	Title       Name       `json:"title"`
	Category    PhoneClass `json:"category"`
	Description PhoneClass `json:"description"`
	Summary     PhoneClass `json:"summary"`
	Role        PhoneClass `json:"role"`
	URL         EmailClass `json:"url"`
	Media       Media      `json:"media"`
	Repo        EmailClass `json:"repo"`
	Start       EmailClass `json:"start"`
	End         EmailClass `json:"end"`
	Highlights  Commitment `json:"highlights"`
	Location    PhoneClass `json:"location"`
	Keywords    Commitment `json:"keywords"`
}

type Media struct {
	Type            CommitmentType `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Description     string         `json:"description"`
	Items           MediaItems     `json:"items"`
}

type MediaItems struct {
	Type                 string            `json:"type"`
	AdditionalProperties bool              `json:"additionalProperties"`
	Properties           CunningProperties `json:"properties"`
}

type CunningProperties struct {
	Category Name       `json:"category"`
	Name     PhoneClass `json:"name"`
	URL      PhoneClass `json:"url"`
}

type Reading struct {
	Type            CommitmentType `json:"type"`
	Description     string         `json:"description"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           ReadingItems   `json:"items"`
}

type ReadingItems struct {
	Type                 string            `json:"type"`
	AdditionalProperties bool              `json:"additionalProperties"`
	Properties           MagentaProperties `json:"properties"`
}

type MagentaProperties struct {
	Title    Name        `json:"title"`
	Category PhoneClass  `json:"category"`
	URL      EmailClass  `json:"url"`
	Author   *Author     `json:"author,omitempty"`
	Date     EmailClass  `json:"date"`
	Summary  PhoneClass  `json:"summary"`
	From     *PhoneClass `json:"from,omitempty"`
}

type Author struct {
	Type            []string   `json:"type"`
	AdditionalItems bool       `json:"additionalItems"`
	Description     string     `json:"description"`
	Items           PhoneClass `json:"items"`
}

type References struct {
	Type            CommitmentType  `json:"type"`
	Description     string          `json:"description"`
	AdditionalItems bool            `json:"additionalItems"`
	Items           ReferencesItems `json:"items"`
}

type ReferencesItems struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           FriskyProperties `json:"properties"`
}

type FriskyProperties struct {
	Name     Name       `json:"name"`
	Role     PhoneClass `json:"role"`
	Category PhoneClass `json:"category"`
	Private  PhoneClass `json:"private"`
	Summary  PhoneClass `json:"summary"`
	Contact  Other      `json:"contact"`
}

type Samples struct {
	Type            CommitmentType `json:"type"`
	Description     string         `json:"description"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           SamplesItems   `json:"items"`
}

type SamplesItems struct {
	Type                 string                `json:"type"`
	AdditionalProperties bool                  `json:"additionalProperties"`
	Properties           MischievousProperties `json:"properties"`
}

type MischievousProperties struct {
	Title      Name       `json:"title"`
	Summary    PhoneClass `json:"summary"`
	URL        EmailClass `json:"url"`
	Date       EmailClass `json:"date"`
	Highlights Commitment `json:"highlights"`
	Keywords   Commitment `json:"keywords"`
}

type Skills struct {
	Type                 string           `json:"type"`
	Description          string           `json:"description"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           SkillsProperties `json:"properties"`
}

type SkillsProperties struct {
	Sets Sets `json:"sets"`
	List List `json:"list"`
}

type List struct {
	Type            CommitmentType `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           ListItems      `json:"items"`
}

type ListItems struct {
	Type                 string                  `json:"type"`
	AdditionalProperties bool                    `json:"additionalProperties"`
	Properties           BraggadociousProperties `json:"properties"`
}

type BraggadociousProperties struct {
	Name    Name       `json:"name"`
	Level   PhoneClass `json:"level"`
	Summary PhoneClass `json:"summary"`
	Years   Years      `json:"years"`
}

type Sets struct {
	Type            CommitmentType `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Optional        bool           `json:"optional"`
	Items           SetsItems      `json:"items"`
}

type SetsItems struct {
	Type                 string      `json:"type"`
	AdditionalProperties bool        `json:"additionalProperties"`
	Properties           Properties1 `json:"properties"`
}

type Properties1 struct {
	Name   Name       `json:"name"`
	Level  PhoneClass `json:"level"`
	Skills Commitment `json:"skills"`
}

type Social struct {
	Type            CommitmentType `json:"type"`
	Description     string         `json:"description"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           SocialItems    `json:"items"`
}

type SocialItems struct {
	Type                 string      `json:"type"`
	AdditionalProperties bool        `json:"additionalProperties"`
	Properties           Properties2 `json:"properties"`
}

type Properties2 struct {
	Network Name       `json:"network"`
	User    Name       `json:"user"`
	URL     Name       `json:"url"`
	Label   PhoneClass `json:"label"`
}

type Speaking struct {
	Type            CommitmentType `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Section         string         `json:"section"`
	Items           SpeakingItems  `json:"items"`
}

type SpeakingItems struct {
	Type                 string      `json:"type"`
	AdditionalProperties bool        `json:"additionalProperties"`
	Properties           Properties3 `json:"properties"`
}

type Properties3 struct {
	Title      PhoneClass `json:"title"`
	Event      Name       `json:"event"`
	Location   PhoneClass `json:"location"`
	Date       EmailClass `json:"date"`
	Highlights Commitment `json:"highlights"`
	Keywords   Commitment `json:"keywords"`
	Summary    PhoneClass `json:"summary"`
}

type Testimonials struct {
	Type            CommitmentType    `json:"type"`
	Description     string            `json:"description"`
	AdditionalItems bool              `json:"additionalItems"`
	Items           TestimonialsItems `json:"items"`
}

type TestimonialsItems struct {
	Type                 string      `json:"type"`
	AdditionalProperties bool        `json:"additionalProperties"`
	Properties           Properties4 `json:"properties"`
}

type Properties4 struct {
	Name     Name       `json:"name"`
	Quote    Name       `json:"quote"`
	Category PhoneClass `json:"category"`
	Private  PhoneClass `json:"private"`
}

type Writing struct {
	Type            CommitmentType `json:"type"`
	Description     string         `json:"description"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           WritingItems   `json:"items"`
}

type WritingItems struct {
	Type                 string      `json:"type"`
	AdditionalProperties bool        `json:"additionalProperties"`
	Properties           Properties5 `json:"properties"`
}

type Properties5 struct {
	Title     Name       `json:"title"`
	Category  PhoneClass `json:"category"`
	Publisher Publisher  `json:"publisher"`
	Date      EmailClass `json:"date"`
	URL       PhoneClass `json:"url"`
	Summary   PhoneClass `json:"summary"`
}

type Publisher struct {
	Type                 []string            `json:"type"`
	Description          string              `json:"description"`
	Optional             bool                `json:"optional"`
	AdditionalProperties bool                `json:"additionalProperties"`
	Properties           PublisherProperties `json:"properties"`
}

type PublisherProperties struct {
	Name PhoneClass `json:"name"`
	URL  PhoneClass `json:"url"`
}

type SummaryType string

const (
	Boolean SummaryType = "boolean"
	Integer SummaryType = "integer"
	String  SummaryType = "string"
)

type Format string

const (
	Date  Format = "date"
	Email Format = "email"
	URI   Format = "uri"
)

type CommitmentType string

const (
	Array CommitmentType = "array"
)

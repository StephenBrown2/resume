package main

type JSONResume struct {
	Schema               string               `json:"$schema"`
	ID                   string               `json:"$id"`
	AdditionalProperties bool                 `json:"additionalProperties"`
	Definitions          Definitions          `json:"definitions"`
	Properties           JSONResumeProperties `json:"properties"`
	Title                string               `json:"title"`
	Type                 string               `json:"type"`
}

type Definitions struct {
	Iso8601 Iso8601 `json:"iso8601"`
}

type Iso8601 struct {
	Type        Type   `json:"type"`
	Description string `json:"description"`
	Pattern     string `json:"pattern"`
}

type JSONResumeProperties struct {
	Schema       Schema       `json:"$schema"`
	Basics       Basics       `json:"basics"`
	Work         Work         `json:"work"`
	Volunteer    Volunteer    `json:"volunteer"`
	Education    Education    `json:"education"`
	Awards       Awards       `json:"awards"`
	Certificates Certificates `json:"certificates"`
	Publications Publications `json:"publications"`
	Skills       Skills       `json:"skills"`
	Languages    Languages    `json:"languages"`
	Interests    Interests    `json:"interests"`
	References   References   `json:"references"`
	Projects     Projects     `json:"projects"`
	Meta         Meta         `json:"meta"`
}

type Awards struct {
	Type            string      `json:"type"`
	Description     string      `json:"description"`
	AdditionalItems bool        `json:"additionalItems"`
	Items           AwardsItems `json:"items"`
}

type AwardsItems struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           PurpleProperties `json:"properties"`
}

type PurpleProperties struct {
	Title   Image `json:"title"`
	Date    Date  `json:"date"`
	Awarder Image `json:"awarder"`
	Summary Image `json:"summary"`
}

type Image struct {
	Type        Type   `json:"type"`
	Description string `json:"description"`
}

type Date struct {
	Ref Ref `json:"$ref"`
}

type Basics struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           BasicsProperties `json:"properties"`
}

type BasicsProperties struct {
	Name     Name     `json:"name"`
	Label    Image    `json:"label"`
	Image    Image    `json:"image"`
	Email    Schema   `json:"email"`
	Phone    Image    `json:"phone"`
	URL      Schema   `json:"url"`
	Summary  Image    `json:"summary"`
	Location Location `json:"location"`
	Profiles Profiles `json:"profiles"`
}

type Schema struct {
	Type        Type   `json:"type"`
	Description string `json:"description"`
	Format      Format `json:"format"`
}

type Location struct {
	Type                 string             `json:"type"`
	AdditionalProperties bool               `json:"additionalProperties"`
	Properties           LocationProperties `json:"properties"`
}

type LocationProperties struct {
	Address     Image `json:"address"`
	PostalCode  Name  `json:"postalCode"`
	City        Name  `json:"city"`
	CountryCode Image `json:"countryCode"`
	Region      Image `json:"region"`
}

type Name struct {
	Type Type `json:"type"`
}

type Profiles struct {
	Type            string        `json:"type"`
	Description     string        `json:"description"`
	AdditionalItems bool          `json:"additionalItems"`
	Items           ProfilesItems `json:"items"`
}

type ProfilesItems struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           FluffyProperties `json:"properties"`
}

type FluffyProperties struct {
	Network  Image  `json:"network"`
	Username Image  `json:"username"`
	URL      Schema `json:"url"`
}

type Certificates struct {
	Type            string            `json:"type"`
	Description     string            `json:"description"`
	AdditionalItems bool              `json:"additionalItems"`
	Items           CertificatesItems `json:"items"`
}

type CertificatesItems struct {
	Type                 string              `json:"type"`
	AdditionalProperties bool                `json:"additionalProperties"`
	Properties           TentacledProperties `json:"properties"`
}

type TentacledProperties struct {
	Name   Image  `json:"name"`
	Date   Date   `json:"date"`
	URL    Schema `json:"url"`
	Issuer Image  `json:"issuer"`
}

type Education struct {
	Type            string         `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           EducationItems `json:"items"`
}

type EducationItems struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           StickyProperties `json:"properties"`
}

type StickyProperties struct {
	Institution Image   `json:"institution"`
	URL         Schema  `json:"url"`
	Area        Image   `json:"area"`
	StudyType   Image   `json:"studyType"`
	StartDate   Date    `json:"startDate"`
	EndDate     Date    `json:"endDate"`
	Score       Image   `json:"score"`
	Courses     Courses `json:"courses"`
}

type Courses struct {
	Type            string  `json:"type"`
	Description     *string `json:"description,omitempty"`
	AdditionalItems bool    `json:"additionalItems"`
	Items           Image   `json:"items"`
}

type Interests struct {
	Type            string         `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           InterestsItems `json:"items"`
}

type InterestsItems struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           IndigoProperties `json:"properties"`
}

type IndigoProperties struct {
	Name     Image   `json:"name"`
	Keywords Courses `json:"keywords"`
}

type Languages struct {
	Type            string         `json:"type"`
	Description     string         `json:"description"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           LanguagesItems `json:"items"`
}

type LanguagesItems struct {
	Type                 string             `json:"type"`
	AdditionalProperties bool               `json:"additionalProperties"`
	Properties           IndecentProperties `json:"properties"`
}

type IndecentProperties struct {
	Language Image `json:"language"`
	Fluency  Image `json:"fluency"`
}

type Meta struct {
	Type                 string         `json:"type"`
	Description          string         `json:"description"`
	AdditionalProperties bool           `json:"additionalProperties"`
	Properties           MetaProperties `json:"properties"`
}

type MetaProperties struct {
	Canonical    Schema `json:"canonical"`
	Version      Image  `json:"version"`
	LastModified Image  `json:"lastModified"`
}

type Projects struct {
	Type            string        `json:"type"`
	Description     string        `json:"description"`
	AdditionalItems bool          `json:"additionalItems"`
	Items           ProjectsItems `json:"items"`
}

type ProjectsItems struct {
	Type                 string              `json:"type"`
	AdditionalProperties bool                `json:"additionalProperties"`
	Properties           HilariousProperties `json:"properties"`
}

type HilariousProperties struct {
	Name        Image   `json:"name"`
	Description Image   `json:"description"`
	Highlights  Courses `json:"highlights"`
	Keywords    Courses `json:"keywords"`
	StartDate   Date    `json:"startDate"`
	EndDate     Date    `json:"endDate"`
	URL         Schema  `json:"url"`
	Roles       Courses `json:"roles"`
	Entity      Image   `json:"entity"`
	Type        Image   `json:"type"`
}

type Publications struct {
	Type            string            `json:"type"`
	Description     string            `json:"description"`
	AdditionalItems bool              `json:"additionalItems"`
	Items           PublicationsItems `json:"items"`
}

type PublicationsItems struct {
	Type                 string              `json:"type"`
	AdditionalProperties bool                `json:"additionalProperties"`
	Properties           AmbitiousProperties `json:"properties"`
}

type AmbitiousProperties struct {
	Name        Image  `json:"name"`
	Publisher   Image  `json:"publisher"`
	ReleaseDate Date   `json:"releaseDate"`
	URL         Schema `json:"url"`
	Summary     Image  `json:"summary"`
}

type References struct {
	Type            string          `json:"type"`
	Description     string          `json:"description"`
	AdditionalItems bool            `json:"additionalItems"`
	Items           ReferencesItems `json:"items"`
}

type ReferencesItems struct {
	Type                 string            `json:"type"`
	AdditionalProperties bool              `json:"additionalProperties"`
	Properties           CunningProperties `json:"properties"`
}

type CunningProperties struct {
	Name      Image `json:"name"`
	Reference Image `json:"reference"`
}

type Skills struct {
	Type            string      `json:"type"`
	Description     string      `json:"description"`
	AdditionalItems bool        `json:"additionalItems"`
	Items           SkillsItems `json:"items"`
}

type SkillsItems struct {
	Type                 string            `json:"type"`
	AdditionalProperties bool              `json:"additionalProperties"`
	Properties           MagentaProperties `json:"properties"`
}

type MagentaProperties struct {
	Name     Image   `json:"name"`
	Level    Image   `json:"level"`
	Keywords Courses `json:"keywords"`
}

type Volunteer struct {
	Type            string         `json:"type"`
	AdditionalItems bool           `json:"additionalItems"`
	Items           VolunteerItems `json:"items"`
}

type VolunteerItems struct {
	Type                 string           `json:"type"`
	AdditionalProperties bool             `json:"additionalProperties"`
	Properties           FriskyProperties `json:"properties"`
}

type FriskyProperties struct {
	Organization Image   `json:"organization"`
	Position     Image   `json:"position"`
	URL          Schema  `json:"url"`
	StartDate    Date    `json:"startDate"`
	EndDate      Date    `json:"endDate"`
	Summary      Image   `json:"summary"`
	Highlights   Courses `json:"highlights"`
}

type Work struct {
	Type            string    `json:"type"`
	AdditionalItems bool      `json:"additionalItems"`
	Items           WorkItems `json:"items"`
}

type WorkItems struct {
	Type                 string                `json:"type"`
	AdditionalProperties bool                  `json:"additionalProperties"`
	Properties           MischievousProperties `json:"properties"`
}

type MischievousProperties struct {
	Name        Image   `json:"name"`
	Location    Image   `json:"location"`
	Description Image   `json:"description"`
	Position    Image   `json:"position"`
	URL         Schema  `json:"url"`
	StartDate   Date    `json:"startDate"`
	EndDate     Date    `json:"endDate"`
	Summary     Image   `json:"summary"`
	Highlights  Courses `json:"highlights"`
}

type Type string

const (
	String Type = "string"
)

type Ref string

const (
	DefinitionsIso8601 Ref = "#/definitions/iso8601"
)

type Format string

const (
	Email Format = "email"
	URI   Format = "uri"
)

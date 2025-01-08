package models

type Radical struct {
	ID       int
	Glyph    string
	Meanings []RadicalMeaning
}

type RadicalMeaning struct {
	ID      *int
	Meaning *string
}

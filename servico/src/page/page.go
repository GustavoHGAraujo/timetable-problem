package page

type Section struct {
  Title         string          `json:"section_title"`
  Content       string        `json:"content"`
}

type Sections []Section

type Page struct {
  Title         string        `json:"page_title"`
  Sections      Sections      `json:"page_sections"`
}

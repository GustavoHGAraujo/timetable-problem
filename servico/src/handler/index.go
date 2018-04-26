package handler

import (
  "html/template"
  "net/http"
  db "../database"
  "../util"
)

type IndexPage struct {
  Title         string           `json:"title"`
  Professores   []db.Professor   `json:"professores"`
  Disciplina    []db.Disciplina  `json:"disciplinas"`
}

func loadIndexPage() (IndexPage) {

  professores, _ := db.GetProfessores()
  disciplinas, _ := db.GetDisciplinas()
  return IndexPage { "Página principal", professores, disciplinas }
}

func Index(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "Index()")

  p := loadIndexPage()

  t := template.Must(template.ParseFiles(
    "pages/index.html",
    "pages/view/professores.html",
    "pages/view/adapter/professor.html",
    "pages/view/disciplinas.html",
    "pages/view/adapter/disciplina.html"))
  t.Execute(w, p)
}

func IndexJSON(w http.ResponseWriter, r *http.Request) {
  // p := criaDummyPageContent()
  // printJson(p, w, r)
}

package handlerView

import (
  "html/template"
  "net/http"
  "log"
  db "../database"
  "../util"
  "../simplex"
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
    "pages/section-professores.html",
    "pages/section-disciplinas.html",
    "pages/item-professor.html",
    "pages/item-disciplina.html"))
  t.Execute(w, p)
}

func IndexJSON(w http.ResponseWriter, r *http.Request) {
  // p := criaDummyPageContent()
  // printJson(p, w, r)
}

func RunSimplex(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "RunSimplex()")
  professores, _ := db.GetProfessores()
  disciplinas, _ := db.GetDisciplinas()
  c := make(chan simplex.SimplexResult)

  go simplex.Run(c, professores, disciplinas)

  result := <- c
  log.Println("Result:", result.Result)
  if result.Error != nil {
    log.Println("Error: ", result.Error.Error())
  }

  if !util.CheckError(result.Error, util.PrintErrorJson, w, r) { util.PrintJson(result.Result, w, r); }
}

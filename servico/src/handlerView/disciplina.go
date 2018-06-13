package handlerView

import (
  "html/template"
  "net/http"
  "strconv"

  "../util"
  "../url"
  db "../database"

  "github.com/gorilla/mux"
)

const TAG = "Handler-View"

func NewDisciplina(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "NewDisciplina()")

  if r.Method == "GET" {
    t := template.Must(template.ParseFiles(
      "pages/new-disciplina.html",
      "pages/option-professor.html"))

    professores, err := db.GetProfessores();
    if !util.CheckError(err, util.ShowError, w, r) { t.Execute(w, professores); }
  } else {
    r.ParseForm()

    disciplina, err := db.ScanDisciplina(r.Form)
    if util.CheckError(err, util.CadastroError, w, r) { return; }

    err = disciplina.New()
    if util.CheckError(err, util.CadastroError, w, r) { return; }

    t, err := template.ParseFiles("pages/info-register-success.html")
    if !util.CheckError(err, util.CadastroError, w, r) { t.Execute(w, nil); }
  }
}

func GetDisciplinaById(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "GetDisciplinaById()")
  r.ParseForm()

  vars := mux.Vars(r)
  disciplinaId, err := strconv.Atoi(vars[url.ID_DISCIPLINA])
  if util.CheckError(err, util.PrintErrorJson, w, r) { return; }

  t, err := template.ParseFiles("pages/view-disciplina.html")
  if util.CheckError(err, util.ShowError, w, r) { return; }

  disciplina, err := db.GetDisciplinaById(disciplinaId)
  if !util.CheckError(err, util.ShowError, w, r) { t.Execute(w, disciplina); }
  return
}

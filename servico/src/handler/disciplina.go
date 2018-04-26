package handler

import (
  "html/template"
  "net/http"
  "strconv"

  "../util"
  "../url"
  db "../database"

  "github.com/gorilla/mux"
)

func NewDisciplina(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "NewDisciplina()")

  if r.Method == "GET" {
    t, err := template.ParseFiles("pages/view/new/disciplina.html")
    if !CheckError(err, ShowError, w, r) { t.Execute(w, nil); }
  } else {
    r.ParseForm()

    disciplina, err := db.ScanDisciplina(r.Form)
    if CheckError(err, CadastroError, w, r) { return; }

    err = disciplina.New()
    if CheckError(err, CadastroError, w, r) { return; }

    t, err := template.ParseFiles("pages/view/msg/cadastro/sucesso.html")
    if !CheckError(err, CadastroError, w, r) { t.Execute(w, nil); }
  }
}

func GetDisciplinaById(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "GetDisciplinaById()")
  r.ParseForm()

  vars := mux.Vars(r)
  disciplinaId, err := strconv.Atoi(vars[url.ID_DISCIPLINA])
  if CheckError(err, PrintErrorJson, w, r) { return; }

  t, err := template.ParseFiles("pages/view/disciplina.html")
  if CheckError(err, ShowError, w, r) { return; }

  disciplina, err := db.GetDisciplinaById(disciplinaId)
  if !CheckError(err, ShowError, w, r) { t.Execute(w, disciplina); }
  return
}

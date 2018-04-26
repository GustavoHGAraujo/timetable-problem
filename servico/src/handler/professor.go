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

func NewProfessor(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "NovoProfessor()")

  if r.Method == "GET" {
    t, err := template.ParseFiles("pages/view/new/professor.html")
    if !CheckError(err, ShowError, w, r) { t.Execute(w, nil); }
  } else {
    r.ParseForm()

    professor, err := db.ScanProfessor(r.Form)
    if CheckError(err, CadastroError, w, r) { return; }

    err = professor.New()
    if CheckError(err, CadastroError, w, r) { return; }

    t, err := template.ParseFiles("pages/view/msg/cadastro/sucesso.html")
    if !CheckError(err, CadastroError, w, r) { t.Execute(w, nil); }
  }
}

func GetProfessorById(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "GetProfessorById()")
  r.ParseForm()

  vars := mux.Vars(r)
  professorId, err := strconv.Atoi(vars[url.ID_PROFESSOR])
  if CheckError(err, PrintErrorJson, w, r) { return; }

  t, err := template.ParseFiles("pages/view/professor.html")
  if CheckError(err, ShowError, w, r) { return; }

  professor, err := db.GetProfessorById(professorId)
  if !CheckError(err, ShowError, w, r) { t.Execute(w, professor); }
  return
}

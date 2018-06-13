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

func NewProfessor(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "NovoProfessor()")

  if r.Method == "GET" {
    t, err := template.ParseFiles("pages/new-professor.html")
    if !util.CheckError(err, util.ShowError, w, r) { t.Execute(w, nil); }
  } else {
    r.ParseForm()

    professor, err := db.ScanProfessor(r.Form)
    if util.CheckError(err, util.CadastroError, w, r) { return; }

    err = professor.New()
    if util.CheckError(err, util.CadastroError, w, r) { return; }

    t, err := template.ParseFiles("pages/info-register-success.html")
    if !util.CheckError(err, util.CadastroError, w, r) { t.Execute(w, nil); }
  }
}

func EditProfessor(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "EditProfessor()")
  r.ParseForm()

  vars := mux.Vars(r)
  professorId, err := strconv.Atoi(vars[url.ID_PROFESSOR])
  if util.CheckError(err, util.PrintErrorJson, w, r) { return; }

  if r.Method == "GET" {
    professor, err := db.GetProfessorById(professorId)
    if util.CheckError(err, util.ShowError, w, r) { return; }

    t, err := template.ParseFiles("pages/edit-professor.html")
    if !util.CheckError(err, util.ShowError, w, r) { t.Execute(w, professor); }
  } else {
    professor, err := db.ScanProfessor(r.Form)
    if util.CheckError(err, util.CadastroError, w, r) { return; }

    professor.Id = professorId
    err = professor.Update()
    if util.CheckError(err, util.CadastroError, w, r) { return; }

    t, err := template.ParseFiles("pages/info-update-success.html")
    if !util.CheckError(err, util.EditarError, w, r) { t.Execute(w, nil); }
  }
  return
}

func DeleteProfessorById(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "DeleteProfessorById()")
  r.ParseForm()

  vars := mux.Vars(r)
  professorId, err := strconv.Atoi(vars[url.ID_PROFESSOR])
  if util.CheckError(err, util.PrintErrorJson, w, r) { return; }

  professor, err := db.GetProfessorById(professorId)
  if util.CheckError(err, util.ShowError, w, r) { return; }

  t, err := template.ParseFiles("pages/view-professor.html")
  if !util.CheckError(err, util.ShowError, w, r) { t.Execute(w, professor); }
  return
}

func GetProfessorById(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "GetProfessorById()")
  r.ParseForm()

  vars := mux.Vars(r)
  professorId, err := strconv.Atoi(vars[url.ID_PROFESSOR])
  if util.CheckError(err, util.PrintErrorJson, w, r) { return; }

  t, err := template.ParseFiles("pages/view-professor.html")
  if util.CheckError(err, util.ShowError, w, r) { return; }

  professor, err := db.GetProfessorById(professorId)
  if !util.CheckError(err, util.ShowError, w, r) { t.Execute(w, professor); }
  return
}

func GetProfessores(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "GetProfessores()")

  t := template.Must(template.ParseFiles(
    "pages/view-professores.html",
    "pages/item-professor.html"))

  professores, err := db.GetProfessores()
  data := struct {Professores db.Professores} {professores}

  if !util.CheckError(err, util.ShowError, w, r) { t.Execute(w, data); }
  return
}

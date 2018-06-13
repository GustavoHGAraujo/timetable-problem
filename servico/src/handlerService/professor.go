package handlerService

import (
  "net/http"
  "strconv"

  "../util"
  "../url"
  db "../database"

  "github.com/gorilla/mux"
)

const TAG = "Handler-Service"

func GetProfessorTodos(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "GetProfessorTodos()")

  professores, err := db.GetProfessores()
  if !util.CheckError(err, util.PrintErrorJson, w, r) { util.PrintJson(professores, w, r); }
  return
}

func GetProfessorById(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "GetProfessorById()")

  vars := mux.Vars(r)
  professorId, err := strconv.Atoi(vars[url.ID_PROFESSOR])
  if util.CheckError(err, util.PrintErrorJson, w, r) { return; }

  professor, err := db.GetProfessorById(professorId)
  if !util.CheckError(err, util.PrintErrorJson, w, r) { util.PrintJson(professor, w, r); }
}

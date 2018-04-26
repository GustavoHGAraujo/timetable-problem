package json

import (
  "net/http"

  handler "../../handler"
  "../../util"
  db "../../database"
)

const TAG = "JSON Handler"

func GetProfessorTodos(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "GetProfessorTodos()")

  professores, err := db.GetProfessores()
  if !handler.CheckError(err, handler.ShowError, w, r) { handler.PrintJson(professores, w, r) }
  return
}

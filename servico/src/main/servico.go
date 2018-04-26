package main

import (
  "log"
  "net/http"
  "os"

  handler "../handler"
  handlerJson "../handler/json"
  "../util"
  "../url"

  "github.com/gorilla/mux"
)

const TAG = "Main"

func main() {
  if len(os.Args) > 1 {
    url.PORT = os.Args[1]
    if len(os.Args) > 2 {
  	   url.DOMAIN = os.Args[2]
     }
  }

	util.LogD("Index URL: " + url.DOMAIN + ":" + url.PORT + url.INDEX)

  router := mux.NewRouter().StrictSlash(true)

  // Resources
  router.HandleFunc("/{FileName}.css", handler.Styles)
  router.HandleFunc("/images/{FileName}.jpg", handler.Images)

  // Páginas
	router.HandleFunc(url.INDEX, handler.Index)

  // CRUD
  router.HandleFunc(url.PROFESSOR_NOVO, handler.NewProfessor)
  router.HandleFunc(url.PROFESSOR, handler.GetProfessorById)
  router.HandleFunc(url.JSON_PROFESSOR_TODOS, handlerJson.GetProfessorTodos)

  router.HandleFunc(url.DISCIPLINA_NOVO, handler.NewDisciplina)

  log.Fatal(http.ListenAndServe(":" + url.PORT, router))
}

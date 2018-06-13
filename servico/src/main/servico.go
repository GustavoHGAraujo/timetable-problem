package main

import (
  "log"
  "net/http"
  "os"

  "../handlerView"
  "../handlerService"
  "../util"
  url "../url"
  urlView "../url/view"
  urlService "../url/service"

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
  router.HandleFunc("/{FileName}.css", handlerView.Styles)
  router.HandleFunc("/images/{FileName}.jpg", handlerView.Images)

  // Páginas
	router.HandleFunc(urlView.INDEX, handlerView.Index)

  // CRUD
  router.HandleFunc(urlView.PROFESSOR_NOVO, handlerView.NewProfessor)
  router.HandleFunc(urlView.PROFESSOR_EDITAR, handlerView.EditProfessor)
  router.HandleFunc(urlView.PROFESSOR, handlerView.GetProfessorById)
  router.HandleFunc(urlView.PROFESSORES, handlerView.GetProfessores)
  router.HandleFunc(urlView.SIMPLEX, handlerView.RunSimplex)
  //router.HandleFunc(urlView.PROFESSOR_DELETE, handlerView.DeleteProfessorById)
  router.HandleFunc(urlService.PROFESSOR_TODOS, handlerService.GetProfessorTodos)

  router.HandleFunc(urlView.DISCIPLINA_NOVO, handlerView.NewDisciplina)
  router.HandleFunc(urlView.DISCIPLINA, handlerView.GetDisciplinaById)

  log.Fatal(http.ListenAndServe(":" + url.PORT, router))
}

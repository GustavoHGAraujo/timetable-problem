package handler

import (
  "net/http"
  "os"
	"time"

  "../util"

	"github.com/gorilla/mux"
)

func Styles(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "Styles()")

  vars := mux.Vars(r)
  fileName := vars["FileName"]
  serve(w, r, "pages/styles/" + fileName + ".css", fileName + ".css")
}

func Images(w http.ResponseWriter, r *http.Request) {
  util.LogD(TAG, "Image()")

  vars := mux.Vars(r)
  fileName := vars["FileName"]
  serve(w, r, "pages/images/" + fileName + ".jpg", fileName + ".jpg")
}

func serve(w http.ResponseWriter, r *http.Request, absolutePath, filename string) {
  file, err := os.Open(absolutePath)
	if err != nil {
		util.LogE(TAG, err);
		PrintJson(err.Error() , w, r)
		return
	}
	defer file.Close()

	util.LogD(TAG, "Serving file.");
	http.ServeContent(w, r, filename, time.Now(), file)
}

package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"../util"
)

const TAG = "Handler"
type ErrorFunction func(err error, w http.ResponseWriter, r *http.Request)
type ErrorPage struct {
	Title		string
	Error		string
}

func CheckError(err error, errorFunction ErrorFunction, w http.ResponseWriter, r *http.Request) (bool) {
	if (err != nil) {
		util.LogE(TAG, err)
		errorFunction(err, w, r)
		return true
	} else {
		return false
	}
}

func toErrorPage(err error, title string) (ErrorPage) {
	return ErrorPage {
		title,
		err.Error(),
	}
}

func PrintJson(content interface {}, w http.ResponseWriter, r *http.Request){
	util.LogD(TAG, "Printing results")
	json, err := json.Marshal(struct { Data interface{} `json:"data"`} {content})
	if err != nil {
		util.LogE(TAG, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func PrintErrorJson(err error, w http.ResponseWriter, r *http.Request) {
	PrintJson(err, w, r)
}

func ShowError(err error, w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("pages/view/error/error.html")
	t.Execute(w, toErrorPage(err, "Algo de errado não está certo"))
}

func CadastroError(err error, w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("pages/view/msg/cadastro/falha.html")
	t.Execute(w, toErrorPage(err, ""))
}

func PageNotDoneYet(err error, w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("pages/view/error/notDoneYet.html")
	t.Execute(w, toErrorPage(err, ""))
}

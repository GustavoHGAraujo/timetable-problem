package database

//import "database/sql"
import "net/url"
import "../util"
import _ "github.com/go-sql-driver/mysql"
import "github.com/jmoiron/sqlx"
import "github.com/go-playground/form"

const TAG = "Database"
const DATABASE = `otimizacao`
const USER = `root`
const PASS = `root`

const PROFESSOR_TABLE = `Professor`
const DISCIPLINA_TABLE = `Disciplina`

var decoder *form.Decoder
var dbProfessores Professores
var dbDisciplinas Disciplinas

func initDatabaseConnection() (*sqlx.DB, error) {
	util.LogD(TAG, "Initializing Database Connection")
	return sqlx.Connect("mysql", USER + ":" + PASS + "@/" + DATABASE + "?parseTime=true")
}

func ScanProfessor(values url.Values) (Professor, error) {
  util.LogD(TAG, "ScanProfessor()")

  var professor Professor
  decoder = form.NewDecoder()
  err := decoder.Decode(&professor, values)
  if err != nil {
    util.LogE(TAG, err)
    return Professor{}, err
  }

	return professor, err
}

func ScanDisciplina(values url.Values) (Disciplina, error) {
  util.LogD(TAG, "ScanDisciplina()")

	type TempDisciplina struct {
		Nome string `form:"nome_disciplina"`
		IdProfessor int `form:"id_professor"`
	}

  var temp TempDisciplina
  decoder = form.NewDecoder()
  err := decoder.Decode(&temp, values)
  if err != nil {
    util.LogE(TAG, err)
    return Disciplina{}, err
  }

	professor, err := GetProfessorById(temp.IdProfessor)
	if err != nil {
    util.LogE(TAG, err)
    return Disciplina{}, err
  }

	return Disciplina{0, temp.Nome, professor}, nil
}

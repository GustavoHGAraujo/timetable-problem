package database

import (
  "../util"
  "log"
)

type Disciplina struct {
  Id          int           `json:"id" db:"id_disciplina"`
  Nome        string        `json:"nome" db:"nome_disciplina" form:"nome_disciplina"`
  Professor                 `json:"professor"`
}

type Disciplinas []Disciplina

const INSERT_DISCIPLINA                = `INSERT INTO Disciplina (nome_disciplina, id_professor) VALUES (:nome_disciplina, :id_professor)`
const SELECT_DISCIPLINA_ALL            = `SELECT * FROM Disciplina NATURAL JOIN Professor;`
const SELECT_DISCIPLINA_BY_ID          = `SELECT * FROM Disciplina NATURAL JOIN Professor WHERE id_disciplina=?;`

func (d *Disciplina) New() (error) {
  util.LogD(TAG, "Disciplina.New()")

  db, err := initDatabaseConnection()
	if err != nil { return err; }
	defer db.Close()

  tx := db.MustBegin()
  log.Println("Disciplina: ", d)
  result, err := tx.NamedExec(INSERT_DISCIPLINA, &d)
  if err != nil { return err; }

  lastId, err := result.LastInsertId()
  if err != nil { return err; }

  d.Id = int(lastId)
  tx.Commit()


  return nil
}

func GetDisciplinaById(id int) (Disciplina, error) {
  util.LogD(TAG, "GetDisciplinaById()")

  db, err := initDatabaseConnection()
	if err != nil { return Disciplina{}, err; }
	defer db.Close()

  var disciplina Disciplina
  err = db.Get(&disciplina, SELECT_DISCIPLINA_BY_ID, id)
  if err != nil { return Disciplina{}, err; }

  return disciplina, nil
}

func GetDisciplinasById(ids []int) (Disciplinas, error){
  util.LogD(TAG, "GetDisciplinasById()")

  disciplinas := make([]Disciplina, len(ids))
  var err error

  for i, id := range ids {
    disciplinas[i], err = GetDisciplinaById(id)

    if err != nil {
      return Disciplinas{}, err
    }
  }

  return disciplinas, nil
}

func GetDisciplinas() (Disciplinas, error) {
  util.LogD(TAG, "GetDisciplinas()")

  db, err := initDatabaseConnection()
	if err != nil { return Disciplinas{}, err; }
	defer db.Close()

  var disciplinas []Disciplina
  err = db.Select(&disciplinas, SELECT_DISCIPLINA_ALL)
  if err != nil { return Disciplinas{}, err; }

  return disciplinas, nil
}

func (d *Disciplinas) Count() int {
  return len(*d)
}
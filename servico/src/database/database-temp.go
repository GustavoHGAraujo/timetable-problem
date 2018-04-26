package database

import "errors"
import "../util"

type Database struct {
  professores       Professores
  disciplinas       Disciplinas

  lastProfessorId   int
  lastDisciplinaId  int
}

var tempDB = Database{}
var ERROR_DUPLICATE = errors.New("Esse registro já existe no banco de dados")
var ERROR_PROFESSOR_NOT_FOUND = errors.New("Professor não existente no banco de dados.")
var ERROR_DISCIPLINA_NOT_FOUND = errors.New("Disciplina não existente no banco de dados.")

func (db *Database) NewProfessor(p *Professor) (error) {
  util.LogD(TAG, "Database.NewProfessor()")
  for _, p1 := range db.professores {
    if p.Id == p1.Id || p.Nome == p1.Nome {
      return ERROR_DUPLICATE
    }
  }

  db.lastProfessorId++
  p.Id = db.lastProfessorId
  db.professores = append(db.professores, *p)

  return nil
}

func (db *Database) NewDisciplina(d *Disciplina) (error) {
  util.LogD(TAG, "Database.NewDisciplina()")
  for _, disciplina := range db.disciplinas {
    if d.Id == disciplina.Id || d.Nome == disciplina.Nome {
      return ERROR_DUPLICATE
    }
  }

  db.lastDisciplinaId++
  d.Id = db.lastDisciplinaId
  db.disciplinas = append(db.disciplinas, *d)

  return nil
}

func (db *Database) GetProfessorById(id int) (Professor, error) {
  util.LogD(TAG, "Database.GetProfessorById()")
  for _, p := range db.professores {
    if p.Id == id {
      return p, nil
    }
  }

  return Professor{}, ERROR_PROFESSOR_NOT_FOUND
}

func (db *Database) GetDisciplinaById(id int) (Disciplina, error) {
  util.LogD(TAG, "Database.GetDisciplinaById()")
  for _, d := range db.disciplinas {
    if d.Id == id {
      return d, nil
    }
  }

  return Disciplina{}, ERROR_DISCIPLINA_NOT_FOUND
}

func (db *Database) GetProfessores() (Professores, error) {
  util.LogD(TAG, "Database.GetProfessores()")
  return db.professores, nil
}

func (db *Database) GetDisciplinas() (Disciplinas, error) {
  util.LogD(TAG, "Database.GetDisciplinas()")
  return db.disciplinas, nil
}

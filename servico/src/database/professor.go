package database

import (
	// "strconv"
	// "log"
	"../util"
	"github.com/jmoiron/sqlx"
)

type Professor struct {
	Id         int      `json:"id" db:"id_professor"`
	Nome       string   `json:"nome" db:"nome_professor" form:"nome_professor"`
	Prioridade int      `json:"prioridade" db:"prioridade" form:"prioridade"`
	Horarios   [][]bool `json:"disponibilidade" form:"horario"`
}

type Horario struct {
	IdProfessor int `db:"id_professor"`
	DiaSemana   int `db:"dia_semana"`
	Posicao     int `db:"posicao"`
}

type Professores []Professor

const INSERT_PROFESSOR = `INSERT INTO Professor (nome_professor) VALUES (:nome_professor)`
const UPDATE_PROFESSOR = `UPDATE Professor SET nome_professor = :nome_professor WHERE id_professor = :id_professor`
const SELECT_PROFESSOR_ALL = `SELECT * FROM Professor;`
const SELECT_PROFESSOR_BY_ID = `SELECT * FROM Professor WHERE id_professor=?;`

//const UPDATE_PROFESSOR                = `UPDATE Professor SET nome_professor = :nome_professor`
const INSERT_HORARIO = `INSERT INTO Horario (id_professor, dia_semana, posicao) VALUES (:id_professor, :dia_semana, :posicao)`
const REMOVE_HORARIOS = `DELETE FROM Horario WHERE id_professor = :id_professor`
const SELECT_HORARIOS_BY_PROFESSOR_ID = `SELECT * FROM Horario WHERE id_professor=?;`

func (p *Professor) New() (error) {
	util.LogD(TAG, "Professor.New()")

	db, err := initDatabaseConnection()
	if err != nil {
		return err;
	}
	defer db.Close()

	tx := db.MustBegin()
	result, err := tx.NamedExec(INSERT_PROFESSOR, &p)
	if err != nil {
		return err;
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return err;
	}

	p.Id = int(lastId)

	err = p.insertHorariosDisponiveis(tx)
	if err != nil {
		return err;
	}

	tx.Commit()

	return nil
}

func (p *Professor) Update() (error) {
	util.LogD(TAG, "Professor.Update()")

	db, err := initDatabaseConnection()
	if err != nil {
		return err;
	}
	defer db.Close()

	tx := db.MustBegin()
	_, err = tx.NamedExec(UPDATE_PROFESSOR, &p)
	if err != nil {
		return err;
	}

	err = p.insertHorariosDisponiveis(tx)
	if err != nil {
		return err;
	}

	tx.Commit()

	return nil
}

func GetProfessorById(id int) (Professor, error) {
	util.LogD(TAG, "GetProfessorById()")

	db, err := initDatabaseConnection()
	if err != nil {
		return Professor{}, err;
	}
	defer db.Close()

	var professor Professor
	err = db.Get(&professor, SELECT_PROFESSOR_BY_ID, id)
	if err != nil {
		return Professor{}, err;
	}

	professor.mergeHorariosDisponiveis()

	return professor, nil
}

func (p *Professor) insertHorariosDisponiveis(tx *sqlx.Tx) (error) {
	_, err := tx.NamedExec(REMOVE_HORARIOS, &p)

	for i, array := range p.Horarios {
		for j, val := range array {
			if val {
				horario := Horario{p.Id, i, j}
				_, err = tx.NamedExec(INSERT_HORARIO, &horario)
				if err != nil {
					return err;
				}
			}
		}
	}

	return nil
}

func (p *Professor) mergeHorariosDisponiveis() (error) {
	db, err := initDatabaseConnection()
	if err != nil {
		return err;
	}
	defer db.Close()

	var horarios []Horario
	err = db.Select(&horarios, SELECT_HORARIOS_BY_PROFESSOR_ID, p.Id)
	if err != nil {
		return err;
	}

	bHorario := [][]bool{{false, false, false},
						{false, false, false},
						{false, false, false},
						{false, false, false},
						{false, false, false}};
	for _, h := range horarios {
		bHorario[h.DiaSemana][h.Posicao] = true
	}

	p.Horarios = bHorario
	return nil
}

func GetProfessoresById(ids []int) (Professores, error) {
	util.LogD(TAG, "GetProfessoresById()")

	professores := make([]Professor, len(ids))
	var err error

	for i, id := range ids {
		professores[i], err = GetProfessorById(id)

		if err != nil {
			return Professores{}, err
		}
	}

	return professores, nil
}

func GetProfessores() (Professores, error) {
	util.LogD(TAG, "GetProfessores()")

	db, err := initDatabaseConnection()
	if err != nil {
		return Professores{}, err;
	}
	defer db.Close()

	var professores []Professor
	err = db.Select(&professores, SELECT_PROFESSOR_ALL)
	if err != nil {
		return Professores{}, err;
	}

	for i := range professores {
		err = professores[i].mergeHorariosDisponiveis()
		if err != nil {
			return Professores{}, err;
		}
	}

	return professores, nil
}

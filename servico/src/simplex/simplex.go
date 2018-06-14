package simplex

import db "../database"

import (
	"fmt"
	"../glpk"
	"../util"
)

const TAG = "Simplex"

type Result struct {
	Board [][]db.Disciplina
	Error error
}

func Run(c chan Result, professores db.Professores, disciplinas db.Disciplinas) {
	util.LogD(TAG, "Run()")

	fmtRestrictionName := "R(3-%02d)(%02d)"
	fmtVariableName := "X(t%02d, s%02d, c%02d)"

	m := New(professores, disciplinas)

	//GLPK
	lp := glpk.New()
	lp.SetProbName("TimeTable")
	lp.SetObjName("Z")
	lp.SetObjDir(glpk.MAX)

	// Creation of the restrictions
	lp.AddRows(m.CountRestrictions())
	rBegin := 1
	rEnd := m.CountClasses()

	// Setting up set 3-1 of restrictions
	for c := rBegin; c <= rEnd; c++ {
		name := fmt.Sprintf(fmtRestrictionName, 1, c)
		lp.SetRowName(c, name)
		lp.SetRowBnds(c, glpk.DB, 0, 1) // <= 1
	}

	// Setting up set 3-2.1 of restrictions
	rBegin = rEnd + 1
	rEnd += m.CountDisciplines()
	for s := rBegin; s <= rEnd; s++ {
		name := fmt.Sprintf(fmtRestrictionName, 2, s)
		lp.SetRowName(s, name)
		lp.SetRowBnds(s, glpk.FX, 2, 2) // == 2
	}

	// Setting up set 3-2.2 of restrictions
	rBegin = rEnd + 1
	rEnd += m.CountDisciplines()
	for s := rBegin; s <= rEnd; s++ {
		name := fmt.Sprintf(fmtRestrictionName, -1, s)
		lp.SetRowName(s, name)
		lp.SetRowBnds(s, glpk.FX, 0, 0) // == 0
	}

	// Setting up set 3-4 of restrictions
	rBegin = rEnd + 1
	rEnd += m.CountTeachers() * m.CountDisciplines() * 5 // 5 dias na semana
	for tc := rBegin; tc <= rEnd; tc++ {
		name := fmt.Sprintf(fmtRestrictionName, 3, tc)
		lp.SetRowName(tc, name)
		lp.SetRowBnds(tc, glpk.DB, 0, 1) // <= 1
	}

	// Setting up set 3-3 of restrictions
	rBegin = rEnd + 1
	rEnd = m.CountRestrictions()
	for tc := rBegin; tc <= rEnd; tc++ {
		name := fmt.Sprintf(fmtRestrictionName, 3, tc)
		lp.SetRowName(tc, name)
		lp.SetRowBnds(tc, glpk.FX, 0, 0) // == 0
	}

	// Creation and setup of the variables
	lp.AddCols(m.CountVariables())
	for t := 0; t < m.CountTeachers(); t++ {
		for c := 0; c < m.CountClasses(); c++ {
			for s := 0; s < m.CountDisciplines(); s++ {
				name := fmt.Sprintf(fmtVariableName, t, s, c)
				l := m.GetVariable(t, s, c)
				lp.SetColName(l, name)
				lp.SetColBnds(l, glpk.DB, 0, 1)
				lp.SetColKind(l, glpk.BV)
			}
		}
	}

	// Adding the coefficients of each variables
	for v := 1; v <= m.CountVariables(); v++ {
		// TODO: Replace coefficient with the teacher preference value
		coefficient := 1.0
		lp.SetObjCoef(v, coefficient)
	}

	// Adding the coefficients of each restriction
	for p := 1; p <= m.CountRestrictions(); p++ {
		lp.SetMatRow(p, m.Variables, m.GetCoefficient(p-1))
	}

	// Run simplex
	err := lp.Simplex(nil)

	// Acquire results
	output := m.GetResult(lp)

	// Erase the problem instance from memory
	lp.Erase()
	lp.Delete()

	// Return the results
	result := new(Result)
	result.Board = output
	result.Error = err

	c <- *result
	util.LogD(TAG, "Run() ended.")
}

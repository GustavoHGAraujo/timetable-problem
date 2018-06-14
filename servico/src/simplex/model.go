package simplex

import (
	"fmt"
	"../glpk"
	"../database"
)

type Model struct {
	Teachers                database.Professores
	Subjects                database.Disciplinas
	RestrictionCoefficients [][]float64
	Variables               []int32
}

func New(teachers database.Professores, disciplines database.Disciplinas) Model {
	model := Model{
		Teachers: teachers,
		Subjects: disciplines,
	}

	model.createVariables()
	model.createRestrictionCoefficients()

	return model
}

func (m *Model) newRestriction() int {
	r := len(m.RestrictionCoefficients)
	m.RestrictionCoefficients = append(m.RestrictionCoefficients, make([]float64, m.CountVariables()))
	for c := 0; c <= m.CountVariables(); c ++ {
		m.RestrictionCoefficients[r][c] = 0.0
	}

	return r
}

func (m *Model) createVariables() {
	numberOfVariables := m.CountClasses() * m.CountTeachers() * m.CountDisciplines()
	variables := make([]int32, numberOfVariables + 1)
	for i := 1; i <= numberOfVariables; i++ {
		variables[i] = int32(i)
	}

	m.Variables = variables
}

func (m *Model) createRestrictionCoefficients() {
	m.RestrictionCoefficients = make([][]float64, 0)

	// Restrictions 3-1
	for k := 0; k < m.CountClasses(); k++ {
		r := m.newRestriction()
		for i := range m.Teachers {
			for j := range m.Subjects {
				l := m.GetVariable(i, j, k)
				m.RestrictionCoefficients[r][l] = 1.0
			}
		}
	}

	// Restrictions 3-2.1
	for j, subject := range m.Subjects {
		r := m.newRestriction()
		i := -1
		for p, teacher := range m.Teachers {
			if teacher.Id == subject.Professor.Id {
				i = p
				break
			}
		}

		if i != -1 {
			for k := 0; k < m.CountClasses(); k++ {
				l := m.GetVariable(i, j, k)
				m.RestrictionCoefficients[r][l] = 1.0
			}
		}
	}

	// Restriction 3-2.2
	for j, subject := range m.Subjects {
		r := m.newRestriction()
		for i, teacher := range m.Teachers {
			for k := 0; k < m.CountClasses(); k++ {
				l := m.GetVariable(i, j, k)
				if teacher.Id != subject.Professor.Id {
					m.RestrictionCoefficients[r][l] = 1.0
				}
			}
		}
	}

	// Restriction 3-4
	for i, teacher := range m.Teachers {
		for j := range m.Subjects {
			for h1, classes := range teacher.Horarios {
				r := m.newRestriction()
				for h2 := range classes {
					k := h1*3 + h2
					l := m.GetVariable(i, j, k)
					m.RestrictionCoefficients[r][l] = 1.0
				}
			}
		}
	}

	// Restriction 3-3
	for i, teacher := range m.Teachers {
		for k := 0; k < m.CountClasses(); k++ {
			h1 := k / 3
			h2 := k % 3
			if !teacher.Horarios[h1][h2] {
				r := m.newRestriction()
				for j := range m.Subjects {
					l := m.GetVariable(i, j, k)
					m.RestrictionCoefficients[r][l] = 1.0
				}
			}
		}
	}
}

func (m *Model) CountClasses() int {
	return 15
}

func (m *Model) CountTeachers() int {
	return m.Teachers.Count()
}

func (m *Model) CountDisciplines() int {
	return m.Subjects.Count()
}

func (m *Model) CountVariables() int {
	return len(m.Variables)
}

func (m *Model) CountRestrictions() int {
	return len(m.RestrictionCoefficients)
}

func (m *Model) GetVariable(p, d, h int) int {
	return h + (d * m.CountClasses()) + (p * m.CountClasses() * m.CountDisciplines()) + 1
}

func (m *Model) GetCoefficient(restriction int) []float64 {
	return m.RestrictionCoefficients[restriction]
}

func (m *Model) GetResult(lp *glpk.Prob) [][]database.Disciplina {
	// Shows results in the console
	output := fmt.Sprintf("Simplex result:\n%s = %g\nVariables:\n", lp.ObjName(), lp.ObjVal())
	for i := 1; i <= m.CountVariables(); i++ {
		output += fmt.Sprintf("%s = %g\n", lp.ColName(i), lp.ColPrim(i))
	}
	fmt.Println(output)

	// Creating subjects matrix
	board := make([][]database.Disciplina, 5)
	for i := range board {
		board[i] = make([]database.Disciplina, 3)
	}

	for l := 0; l < m.CountVariables(); l++ {
		// p := l / (CountClasses * CountDisciplines) // teacher index
		d := (l / m.CountClasses()) % m.CountDisciplines()
		h := l % m.CountClasses()
		h1 := h / 3
		h2 := h % 3

		if lp.ColPrim(l+1) > 0.0 {
			board[h1][h2] = m.Subjects[d]
		}
	}

	return board
}

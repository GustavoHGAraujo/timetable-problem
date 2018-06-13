package simplex

import db "../database"

import (
	"log"
	"fmt"
	"../glpk"
	"../util"
)

type SimplexResult struct {
	Result     [][]db.Disciplina
	Error      error
}

const TAG = "Simplex"

func Run(c chan SimplexResult, professores db.Professores, disciplinas db.Disciplinas) {
	util.LogD(TAG, "Run()")

	fmtNomeRestricao := "R(3-%d)(%d)"
	fmtNomeVariavel := "X(p%d, d%d, h%d)"

	// Variáveis auxiliares
	numHorarios := 15
	numProfessores := len(professores)
	numDisciplinas := len(disciplinas)

	//GLPK
	lp := glpk.New()
	lp.SetProbName("TimeTable")
	lp.SetObjName("Z")
	lp.SetObjDir(glpk.MAX)

	/* Seção de Restrições */
	coefRestricoes := getCoeficientesRestricoes(professores, disciplinas)
	numRestricoes := len(coefRestricoes)

	// Adição das restrições
	log.Println("Adição de Restrições")
	log.Println("numRestrições:", numRestricoes)
	lp.AddRows(numRestricoes)
	r := 0
	rInicial := 1
	rFinal := numHorarios

	// (3-1)
	log.Println("  3-1: rInicial =", rInicial, "rFinal =", rFinal)
	for h := rInicial; h <= rFinal; h++ {
		nomeRestricao := fmt.Sprintf(fmtNomeRestricao, 1, h)
		lp.SetRowName(h, nomeRestricao)
		lp.SetRowBnds(h, glpk.DB, 0, 1) // somatório <= 1
		r++
	}

	// (3-2)
	rInicial = rFinal + 1
	rFinal += numDisciplinas
	log.Println("  3-2: rInicial =", rInicial, "rFinal =", rFinal)
	for d := rInicial; d <= rFinal; d++ {
		nomeRestricao := fmt.Sprintf(fmtNomeRestricao, 2, d)
		lp.SetRowName(d, nomeRestricao)
		lp.SetRowBnds(d, glpk.FX, 2, 2) // somatório== 2
		r++
	}

	rInicial = rFinal + 1
	rFinal += numDisciplinas
	log.Println("  3-?: rInicial =", rInicial, "rFinal =", rFinal)
	for d := rInicial; d<= rFinal; d++ {
		nomeRestricao := fmt.Sprintf(fmtNomeRestricao, -1, d)
		lp.SetRowName(d, nomeRestricao)
		lp.SetRowBnds(d, glpk.FX, 0, 0) // somatorio == 0
		r++
	}

	// (3-4)
	rInicial = rFinal + 1
	rFinal += numProfessores * numDisciplinas * 5
	log.Println("  3-4: rInicial =", rInicial, "rFinal =", rFinal)
	for ph := rInicial; ph <= rFinal; ph++ {
		nomeRestricao := fmt.Sprintf(fmtNomeRestricao, 3, ph)
		lp.SetRowName(ph, nomeRestricao)
		lp.SetRowBnds(ph, glpk.DB, 0, 1) // somatório = 1
		r++
	}

	// (3-3)
	rInicial = rFinal + 1
	rFinal = numRestricoes
	log.Println("  3-3: rInicial =", rInicial, "rFinal =", rFinal)
	for ph := rInicial; ph <= rFinal; ph++ {
		nomeRestricao := fmt.Sprintf(fmtNomeRestricao, 3, ph)
		lp.SetRowName(ph, nomeRestricao)
		lp.SetRowBnds(ph, glpk.FX, 0, 0) // somatório== 0
	}

	/* Seção de Variáveis */
	log.Println("numVariáveis:", numHorarios * numProfessores * numDisciplinas)
	numVariaveis := numHorarios * numProfessores * numDisciplinas

	log.Println("Declaração de variáveis")
	variaveis := make([]int32, numVariaveis+1)
	for i := 1; i <= numVariaveis; i++ {
		variaveis[i] = int32(i)
	}

	// Adição das Variáveis
	log.Println("Adição de Variáveis")
	lp.AddCols(numVariaveis)
	log.Println("  Inicializando variáveis")
	for p := 0; p < numProfessores; p++ {
		for h := 0; h < numHorarios; h++ {
			for d := 0; d < numDisciplinas; d++ {
				nomeVariavel := fmt.Sprintf(fmtNomeVariavel, p, d, h)
				l := getPosicaoVariavel(p, d, h, numHorarios, numDisciplinas, numProfessores)
//				log.Println("    Variavel", l, "\tnome:", nomeVariavel)
				lp.SetColName(l, nomeVariavel)
				lp.SetColBnds(l, glpk.DB, 0, 1)
				lp.SetColKind(l, glpk.BV)
			}
		}
	}

	// Adição dos Coeficientes das variáveis
	for v := 1; v <= numVariaveis; v++ {
		// Coef: Professor.Horario.Preferência
		lp.SetObjCoef(v, 1.0)
	}

	// Aplica os coeficientes às restrições
	log.Println("Aplica os coeficientes às restrições")
	for p := 1; p <= numRestricoes; p++ {
		lp.SetMatRow(p, variaveis, coefRestricoes[p-1])
	}

	// Dump
	if true {
		fmt.Println("Primeiro dump")
		for r := 0; r < numRestricoes; r++ {
			for v := 1; v <= numVariaveis; v++ {
				fmt.Printf("%.0f ", coefRestricoes[r][v])
			}
			fmt.Println()
		}
		fmt.Println("Fim do primeiro Dump")
	}

	if true {
		for p := 1; p <= numRestricoes; p++{
			output := "      "
			ind, _ := lp.MatRow(p)
			for v := int32(1); v <= int32(numVariaveis); v++ {
				o := "0 "
				for i := range ind {
					if  v == ind[i] {
						o = "1 "
					}
				}
				output += o
			}
			fmt.Println(output)
		}
	}
	if true {
		output := "F.O.: "
		for v := 1; v <= numVariaveis; v++ {
			if lp.ObjCoef(v) > 0 {
				output += "1 "
			} else {
				output += "0 "
			}
		}
		fmt.Println(output)
	}

	// Roda o simplex
	lp.Simplex(nil)

	// Resultados
	output := getResult(lp, numVariaveis, numRestricoes, numProfessores, numDisciplinas, numHorarios, disciplinas, professores)

	lp.Delete()

	// Retornando resultados
	rcInstance := new(SimplexResult)
	rcInstance.Result = output
	rcInstance.Error = nil

	c <- *rcInstance
	util.LogD(TAG, "Run() ended.")
}

func newCoefRestricao(coefRestricoes *[][]float64, numVariaveis int) int {
	log.Println("newCoefRestricao()")
	r := len(*coefRestricoes)
	*coefRestricoes = append(*coefRestricoes, make([]float64, numVariaveis + 1))
//	log.Println("  inicializando novo array")
	for c := 0; c <= numVariaveis; c ++ {
//		log.Println("    r:", r, "\tc:", c)
		(*coefRestricoes)[r][c] = 0.0
	}

	return r
}

func getPosicaoVariavel(p, d, h, numHorarios, numDisciplinas, numProfessores int) int {
	return h + (d * numHorarios) + (p * numHorarios * numDisciplinas) + 1
}

func getCoeficientesRestricoes(professores db.Professores, disciplinas db.Disciplinas) [][]float64 {
	log.Println("getCoeficienteRestricoes()")
	// Variáveis auxiliares
	numHorarios := 15
	numProfessores := len(professores)
	numDisciplinas := len(disciplinas)

	numVariaveis := numHorarios * numProfessores * numDisciplinas
	coefRestricoes := make([][]float64, 0)

	// Restrição (3-1)
	log.Println("Restrição 3-1")
	for k := 0; k < numHorarios; k++ {
		r := newCoefRestricao(&coefRestricoes, numVariaveis)
		//log.Println("R ",r,": ")
		for i := range professores {
			for j := range disciplinas {
				l := getPosicaoVariavel(i, j, k, numHorarios, numDisciplinas, numProfessores)
				//log.Println("(",i,",",j,",",k,"):[",l,"]  =  1")
				coefRestricoes[r][l] = 1.0
			}
		}
	}

	// Restrição (3-2)
	log.Println("Restrição 3")
	for j, disciplina := range disciplinas {
		r := newCoefRestricao(&coefRestricoes, numVariaveis)
		i := -1
		log.Println("Restricao 3-2: j:", j)
		for p, professor := range professores {
			if professor.Id == disciplina.Professor.Id {
				i = p; break;
			}
		}
		log.Println("Restricaoo 3-2: i:", i)
		if i != -1 {
			for k := 0; k < numHorarios; k++ {
				l := getPosicaoVariavel(i, j, k, numHorarios, numDisciplinas, numProfessores)
				//log.Println("(",i,",",j,",",k,"):[",l,"]  =  1")
				coefRestricoes[r][l] = 1.0
			}
		}
	}

	// Restricao (3-?)
	log.Println("Restricao 3-?")
	for j, disciplina := range disciplinas {
		r := newCoefRestricao(&coefRestricoes, numVariaveis)
		for i, professor := range professores {
			for k := 0; k < numHorarios; k++ {
				l := getPosicaoVariavel(i, j, k, numHorarios, numDisciplinas, numProfessores)
				if professor.Id != disciplina.Professor.Id {
					coefRestricoes[r][l] = 1.0
				}
			}
		}
	}

	// Restrição (3-4)
	log.Println("Restrição 3-4")
	for i, professor := range professores {
		for j := range disciplinas {
			for h1, dia := range professor.Horarios {
				r := newCoefRestricao(&coefRestricoes, numVariaveis)
				//log.Println("R ",r,": ")

				for h2 := range dia {
					k := h1 * 3 + h2
					l := getPosicaoVariavel(i, j, k, numHorarios, numDisciplinas, numProfessores)
					//log.Println("(",i,",",j,",",k,"):[",l,"]  =  1")

					coefRestricoes[r][l] = 1.0
				}
			}
		}
	}

	// Restrição (3-3)
	log.Println("Restrição 3-3")
	for i, professor := range professores {

		for k := 0; k < numHorarios; k++ {
			h1 := k / 3
			h2 := k % 3

			if !professor.Horarios[h1][h2] {
				r := newCoefRestricao(&coefRestricoes, numVariaveis)
				// log.Println("R ",r,": ")
				// h2 := k % 3

				for j := range disciplinas {
					l := getPosicaoVariavel(i, j, k, numHorarios, numDisciplinas, numProfessores)
					//log.Println("(",i,",",j,",",k,"):[",l,"]  =  1")
					coefRestricoes[r][l] = 1.0
				}
			}
		}
	}

	return coefRestricoes
}

func getResult(lp *glpk.Prob, numVariaveis, numRestricoes, numProfessores, numDisciplinas, numHorarios int, disciplinas db.Disciplinas, professores db.Professores) ([][]db.Disciplina) {
	log.Println("getResult()")

	// Mostrar resultado no console
	output := fmt.Sprintf("Resultado:\n%s = %g\nVariáveis:\n", lp.ObjName(), lp.ObjVal())
	for i := 0; i < numVariaveis; i++ {
		output += fmt.Sprintf("%s = %g\n", lp.ColName(i+1), lp.ColPrim(i+1))
	}
	output += "\n"

	fmt.Println(output)


	// Montando quadro que será retornada
	quadro := make([][]db.Disciplina, 5)
	for i := range quadro {
		quadro[i] = make([]db.Disciplina, 3)
	}

	for l := 0; l < numVariaveis; l++ {
		p := l / (numHorarios * numDisciplinas)
		d := (l / numHorarios) % numDisciplinas
		h := l % numHorarios
		h1 := h / 3
		h2 := h % 3

		sLog := fmt.Sprintf("Percorrendo variáveis:l : %2d, p: %2d, d: %2d, h: %2d -> (h1: %2d, h2: %2d)", l, p, d, h, h1, h2)

		if lp.ColPrim(l + 1) > 0.0 {
			sLog += fmt.Sprintf("\tTEM AULA")
			quadro[h1][h2] = disciplinas[d]

			log.Println(sLog)
		}
	}

	return quadro
}

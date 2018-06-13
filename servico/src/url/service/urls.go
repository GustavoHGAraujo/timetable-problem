package service

import url "../"

const INDEX = url.INDEX + `/service`

const PROFESSOR = INDEX + url.Professor + url.VAR_ID_PROFESSOR
const PROFESSOR_TODOS = INDEX + url.Professor + `/todos`
const PROFESSOR_NOVO = INDEX + url.Professor + url.NOVO
const PROFESSOR_EDITAR = PROFESSOR + url.EDITAR

const DISCIPLINA = INDEX + url.Disciplina + url.VAR_ID_DISCIPLINA
const DISCIPLINA_NOVO = INDEX + url.Disciplina + url.NOVO
const DISCIPLINA_EDITAR = DISCIPLINA + url.EDITAR

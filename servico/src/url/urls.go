package url

var DOMAIN = `http://localhost`
var PORT = `9090`

const INDEX = `/otimizacao`
const NOVO = `/novo`
const EDITAR = `/editar`
const REMOVER = `/remover`
const TODOS = `/todos`
const JSON = `/json`

const ID_PROFESSOR = `IdProfessor`
const ID_DISCIPLINA = `IdDisciplina`

const VAR_ID_PROFESSOR  = `/{` + ID_PROFESSOR + `}`
const VAR_ID_DISCIPLINA  = `/{` + ID_DISCIPLINA + `}`

const professor = `/professor`
const disciplina   = `/disciplina`

const PROFESSOR = INDEX + professor + VAR_ID_PROFESSOR
const PROFESSOR_TODOS = INDEX + professor + `/todos`
const PROFESSOR_NOVO = INDEX + professor + NOVO
const PROFESSOR_EDITAR = PROFESSOR + EDITAR

const JSON_PROFESSOR = INDEX + JSON + professor + VAR_ID_PROFESSOR
const JSON_PROFESSOR_TODOS = INDEX + JSON + professor + TODOS

const DISCIPLINA = INDEX + disciplina + VAR_ID_DISCIPLINA
const DISCIPLINA_NOVO = INDEX + disciplina + NOVO
const DISCIPLINA_EDITAR = INDEX + disciplina + EDITAR

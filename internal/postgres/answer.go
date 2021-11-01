package postgres

type Answer string

const (
	AnswerA Answer = "a"
	AnswerB Answer = "b"
	AnswerC Answer = "c"
	AnswerD Answer = "d"
)

func (answer Answer) IsValid() bool {
	switch answer {
	case AnswerA,
		AnswerB,
		AnswerC,
		AnswerD:
		return true
	}
	return false
}

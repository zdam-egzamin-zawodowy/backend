package models

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Answer string

const (
	AnswerA Answer = "a"
	AnswerB Answer = "b"
	AnswerC Answer = "c"
	AnswerD Answer = "d"
)

func (role Answer) IsValid() bool {
	switch role {
	case AnswerA,
		AnswerB,
		AnswerC,
		AnswerD:
		return true
	}
	return false
}

func (answer Answer) String() string {
	return string(answer)
}

func (answer *Answer) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*answer = Answer(strings.ToLower(str))
	if !answer.IsValid() {
		return fmt.Errorf("%s is not a valid Answer", str)
	}
	return nil
}

func (answer Answer) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(answer.String()))
}

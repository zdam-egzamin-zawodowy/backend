package model

import (
	"fmt"
	"github.com/pkg/errors"
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

func (answer Answer) String() string {
	return string(answer)
}

func (answer *Answer) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return errors.New("enums must be strings")
	}

	*answer = Answer(strings.ToLower(str))
	if !answer.IsValid() {
		return errors.Errorf("%s is not a valid Answer", str)
	}
	return nil
}

func (answer Answer) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(answer.String()))
}

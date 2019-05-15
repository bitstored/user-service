package errors

import (
	"fmt"
	"strings"
)

const (
	ErrKindPasswordDoesntComply Kind = iota
	ErrKindInvalidMail
	ErrKindInvalidJWTToken
	ErrKindInvalidUserID
	ErrKindUnableToDecodeToken
	ErrKindAccountLocked
	ErrKindAccountNotFound
	ErrKindTokenExpired
)

type Kind int

type Err struct {
	Msgs []string
	Kind Kind
}

func NewError(kind Kind, msg string) *Err {
	return &Err{Msgs: []string{msg}, Kind: kind}
}

func NewErrorArr(kind Kind, msgs []string) *Err {
	return &Err{Msgs: msgs, Kind: kind}
}

func (e *Err) AddMessage(msg string) {
	e.Msgs = append(e.Msgs, msg)
}

func (e *Err) AddMessageArr(msgs []string) {
	e.Msgs = append(e.Msgs, msgs...)
}

func (e *Err) Message() string {
	return fmt.Sprintf("%s : %s", e.KindMsg(), strings.Join(e.Msgs, "\n"))
}

func (e *Err) Error() error {
	if e == nil || len(e.Msgs) == 0 {
		return nil
	}
	return fmt.Errorf("%s : %s", e.KindMsg(), strings.Join(e.Msgs, "\n"))
}

func (e *Err) KindMsg() string {
	switch e.Kind {
	case ErrKindInvalidMail:
		return "ErrKindInvalidMail"
	case ErrKindPasswordDoesntComply:
		return "ErrKindPasswordDoesntComply"
	case ErrKindInvalidJWTToken:
		return "ErrKindInvalidJWTToken"
	case ErrKindInvalidUserID:
		return "ErrKindInvalidUserID"
	case ErrKindUnableToDecodeToken:
		return "ErrKindUnableToDecodeToken"
	case ErrKindAccountLocked:
		return "ErrKindAccountLocked"
	case ErrKindAccountNotFound:
		return "ErrKindAccountNotFound"
	case ErrKindTokenExpired:
		return "ErrKindTokenExpired"
	default:
		return "ErrKindUnknown"
	}
}

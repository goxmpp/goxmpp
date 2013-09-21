package mechanism

import (
	"github.com/dotdoom/goxmpp/extensions/features/auth"
)

type Plain auth.Mechanism

func NewPlain() *Plain {
	return &Plain{Name: "PLAIN"}
}

func init() {
	auth.Mechanisms.AddElement(NewPlain())
}

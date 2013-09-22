package mechanism

import (
	"bytes"
	"encoding/base64"
	"github.com/dotdoom/goxmpp/extensions/features/auth"
)

type Plain struct {
	auth.MechanismElement
}

func NewPlain() *Plain {
	return &Plain{MechanismElement: auth.MechanismElement{Name: "PLAIN"}}
}

var usernamePasswordSeparator = []byte{0}

func (self *Plain) Process(a *auth.AuthElement) {
	b, _ := base64.StdEncoding.DecodeString(a.Data)
	user_password := bytes.Split(b, usernamePasswordSeparator)
	println("Username:", string(user_password[1]), "password", string(user_password[2]))
}

func init() {
	auth.Mechanisms.AddElement(NewPlain())
}

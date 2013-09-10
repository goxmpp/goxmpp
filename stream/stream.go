package stream

import "github.com/dotdoom/goxmpp"
import "github.com/dotdoom/goxmpp/xep"

var Registrator = xep.NewRegistrator()

type Stream struct {
	XMLName xml.Name
	ID      string
	From    string
	To      string
	Version string
}

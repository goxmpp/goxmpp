package stream

import "github.com/dotdoom/goxmpp"

type Stream struct {
	XMLName xml.Name
	ID      string
	From    string
	To      string
	Version string
}

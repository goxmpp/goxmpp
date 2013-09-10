package stream

import "github.com/dotdoom/goxmpp"
import "github.com/dotdoom/goxmpp/extensions"

var HandlerRegistrator = extensions.NewHandlerRegistrator()

type Stream struct {
	XMLName xml.Name
	ID      string
	From    string
	To      string
	Version string
}

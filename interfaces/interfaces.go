package interfaces

import "github.com/dotdoom/goxmpp"

type InnerElementAdder interface {
	AddInnerElement(InnerElementAdder) bool
}

type ElementHandler interface {
	HandleElement(goxmpp.StreamWrapper)
}

type XMLHandler interface {
	HandleInnerXML(goxmpp.StreamWrapper) []ElementHandler
}

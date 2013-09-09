package xep

import . "github.com/dotdoom/goxmpp/interfaces"

type Registrator map[string]ElementHandler

func (self *Registrator) Register(key string, handler ElementHandler) {
	r := *self
	r[key] = handler
	self = &r
}

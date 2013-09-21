package mechanism

import (
	"github.com/dotdoom/goxmpp/extensions/features/auth"
)

func init() {
	auth.Mechanisms.AddElement(&auth.Mechanism{Name: "PLAIN"})
}

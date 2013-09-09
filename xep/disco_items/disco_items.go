package disco_items

import "enconding/xml"

type DiscoItemsQuery struct {
	// http://xmpp.org/extensions/xep-0030.html
	XMLName xml.Name `xml:"http://jabber.org/protocol/disco#items query"`
}

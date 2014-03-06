package auth

import "encoding/xml"

type MechanismElement struct {
	XMLName xml.Name `xml:"mechanism"`
	Name    string   `xml:",chardata"`
}

func NewMechanismElement(name string) *MechanismElement {
	return &MechanismElement{
		Name: name,
	}
}

package elements

type Element interface{}

type InnerElementsAdder interface {
	AddInnerElement(Element) bool
}

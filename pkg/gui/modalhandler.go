package gui

type ModalHandler interface {
	ModalMessage(msg string)
	ModalYesNo(msg string, yes func())
}

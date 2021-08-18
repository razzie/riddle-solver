package tui

// ModalHandler is an interface that handles modal dialogs
type ModalHandler interface {
	ModalMessage(msg string)
	ModalYesNo(msg string, yes func())
}

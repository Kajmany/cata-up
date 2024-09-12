package ui

type StatusModel struct {
	Common     *Common
	message    string
	statusType statusType
}

type statusType int

const (
	five = 5
)

package constants

import "fmt"

type Status string

// noinspection GoUnusedConst
const (
	StatusCreated    = "Created"
	StatusSuccess    = "Success"
	StatusFailed     = "FailedToDownload"
	StatusInProgress = "InProgress"
)

var ErrCreateTask = fmt.Errorf("such a task has already been created")

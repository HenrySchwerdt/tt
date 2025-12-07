package tterrors

import "errors"

var (
	ErrProjectExists        = errors.New("project already exists")
	ErrProjectDoesNotExist  = errors.New("project does not exist")
	ErrTimeEntryNotClosed   = errors.New("close current time entry before starting a new one")
	ErrNoOpenTimeEntryFound = errors.New("no open time entry found")
	ErrNoProjectPath        = errors.New("need to provide a project path")
)

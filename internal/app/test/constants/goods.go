package constants

import "fmt"

var (
	ErrIncorrect = fmt.Errorf("link is incorrect. url must start with http:// or https://")
	ErrNotNumber = fmt.Errorf("must be a positive number")
	ErrNotFound  = fmt.Errorf("nothing found")
)

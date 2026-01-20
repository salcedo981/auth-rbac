package officesError

import "errors"

var (
	ErrInstiCodeRequired = errors.New("institution code is required")
	ErrNoBranchesFound   = errors.New("no branches found")
	ErrInvalidInput      = errors.New("invalid input")

	ErrBranchCodeRequired = errors.New("branch code is required")
	ErrNoUnitsFound       = errors.New("no units found")
)

package manifold

import "errors"

var (
	errAPITokenRequired   = errors.New("An API Token is required to use this provider")
	errProjectNotFound    = errors.New("Could not find project for label")
	errResourceNotFound   = errors.New("Could not find resource for label")
	errCredentialNotFound = errors.New("Could not find credential")
	errTeamNotFound       = errors.New("Could not find team for label")
)

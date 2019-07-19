package outlook

import (
	"fmt"
	"time"
)

var (
	// ErrNoAccessToken is returned when a query is executed in a session which was either not given a refreshToken or that failed to retrieve and the access token.
	ErrNoAccessToken = fmt.Errorf("no access token for session")
)

// ErrStatusCode an error thrown when a given http call responds with a bad http status
type ErrStatusCode struct {
	Code                   int
	Message                string
	SuggestedRetryDuration time.Duration
}

func (sce *ErrStatusCode) Error() string {
	return fmt.Sprintf(
		"Call to microsoft's graph api failed with a status code: %d. Reason: %s",
		sce.Code,
		sce.Message,
	)
}

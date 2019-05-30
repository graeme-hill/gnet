package uberserver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStartsAndStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	s := StartUberServer(ctx)

	cancel()

	select {
	case errs := <-s.Done():
		for service, err := range errs {
			require.NoError(t, err, service)
		}
	case <-time.After(1 * time.Second):
		t.FailNow()
	}
}

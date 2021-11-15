package subscription

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionCancellations(t *testing.T) {
	cancellations := newSubscriptionCancellations()
	var ctx context.Context

	t.Run("should add a cancellation func to map", func(t *testing.T) {
		require.Equal(t, 0, cancellations.Count())

		ctx = cancellations.AddWithParent("1", context.Background())
		assert.Equal(t, 1, cancellations.Count())
		assert.NotNil(t, ctx)
	})

	t.Run("should execute cancellation from map", func(t *testing.T) {
		require.Equal(t, 1, cancellations.Count())
		ctxTestFunc := func() bool {
			<-ctx.Done()
			return true
		}

		ok := cancellations.Cancel("1")
		assert.Eventually(t, ctxTestFunc, time.Second, 5*time.Millisecond)
		assert.True(t, ok)
		assert.Equal(t, 0, cancellations.Count())
	})
}

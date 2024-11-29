package packets

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPackets(t *testing.T) {
	t.Run("RunTCPServerAndTCPClient", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			<-time.After(time.Second * 5)
			cancel()
		}()

		go func() {
			<-time.After(time.Second * 1)
			err := RunTCPClient(ctx, "127.0.0.1:6943", func(t time.Time, i1, i2, i3, i4 int64) {
				log.Printf("%#+v", []any{t, i1, i2, i3, i4})
			})
			require.NoError(t, err)
		}()

		err := RunTCPServer(ctx, 6943)
		require.NoError(t, err)
	})
}

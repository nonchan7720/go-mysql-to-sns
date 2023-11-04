package logging

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/nonchan7720/go-storage-to-messenger/pkg/mock"
	"github.com/stretchr/testify/require"
)

func testHTTPClient(transport mock.MockRoundTripper) *http.Client {
	client := http.DefaultClient
	client.Transport = transport
	return client
}

func TestRollbar(t *testing.T) {
	called := false
	defer func() {
		require.True(t, called)
	}()
	transport := func(req *http.Request) (*http.Response, error) {
		called = true
		return &http.Response{
			StatusCode: http.StatusOK,
		}, nil
	}
	conf := RollbarConfig{
		LogLevel: "error",
		Token:    "DUMMY",
		Client:   testHTTPClient(transport),
	}
	conf.Init("local", "v1", "test")
	defer conf.Close()
	h := NewRollbarHandler(&conf)
	defer h.Close()
	log := slog.New(h)
	slog.SetDefault(log)
	slog.Error("This is test")
}

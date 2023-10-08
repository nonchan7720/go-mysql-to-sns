package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabase(t *testing.T) {
	conf := Database{
		TLS: nil,
	}
	require.Nil(t, conf.Tls())
	conf = Database{
		TLS: &TLS{},
	}
	require.NotNil(t, conf.Tls())
	conf = Database{
		TLS: &TLS{
			InsecureSkipVerify: true,
			SeverName:          "aaa",
		},
	}
	tls := conf.Tls()
	require.NotNil(t, tls)
	require.True(t, tls.InsecureSkipVerify)
	require.Equal(t, tls.ServerName, "aaa")
}

func TestTLS(t *testing.T) {
	tls := TLS{
		InsecureSkipVerify: false,
		SeverName:          "",
	}
	conf := tls.Config()
	require.False(t, conf.InsecureSkipVerify)
	require.Empty(t, conf.ServerName)
	tls = TLS{
		InsecureSkipVerify: true,
		SeverName:          "aaa",
	}
	conf = tls.Config()
	require.True(t, conf.InsecureSkipVerify)
	require.Equal(t, conf.ServerName, "aaa")
}

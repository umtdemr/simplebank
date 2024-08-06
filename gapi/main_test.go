package gapi

import (
	"github.com/stretchr/testify/require"
	db "github.com/umtdemr/simplebank/db/sqlc"
	"github.com/umtdemr/simplebank/util"
	"github.com/umtdemr/simplebank/worker"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store, distributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, distributor)
	require.NoError(t, err)
	return server

}

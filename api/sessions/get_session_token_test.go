package sessions_test

import (
	"net/http"
	"testing"

	"github.com/keratin/authn-server/api/sessions"
	"github.com/keratin/authn-server/api/test"
	"github.com/keratin/authn-server/lib/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSessionToken(t *testing.T) {
	app := test.App()
	server := test.Server(app, sessions.Routes(app))
	defer server.Close()

	client := route.NewClient(server.URL).Referred(&app.Config.ApplicationDomains[0])

	t.Run("known account", func(t *testing.T) {
		_, err := app.AccountStore.Create("known@keratin.tech", []byte("pwd"))
		require.NoError(t, err)

		res, err := client.Get("/session/token?username=known@keratin.tech")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// TODO: assert go routine?
	})

	t.Run("unknown account", func(t *testing.T) {
		res, err := client.Get("/session/token?username=unknown@keratin.tech")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)
	})
}

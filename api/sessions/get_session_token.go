package sessions

import (
	"net/http"

	"github.com/keratin/authn-server/api"
	"github.com/keratin/authn-server/services"
)

func getSessionToken(app *api.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := app.AccountStore.FindByUsername(r.FormValue("username"))
		if err != nil {
			panic(err)
		}

		// run in the background so that a timing attack can't enumerate usernames
		go func() {
			err := services.PasswordlessTokenSender(app.Config, account)
			if err != nil {
				app.Reporter.ReportRequestError(err, r)
			}
		}()

		w.WriteHeader(http.StatusOK)
	}
}

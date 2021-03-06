package services

import (
	"github.com/keratin/authn-server/config"
	"github.com/keratin/authn-server/data"
	"github.com/keratin/authn-server/lib/route"
	"github.com/keratin/authn-server/models"
	"github.com/keratin/authn-server/ops"
	"github.com/keratin/authn-server/tokens/identities"
	"github.com/keratin/authn-server/tokens/sessions"
	"github.com/pkg/errors"
)

func SessionCreator(
	accountStore data.AccountStore, refreshTokenStore data.RefreshTokenStore, keyStore data.KeyStore, actives data.Actives, cfg *config.Config, reporter ops.ErrorReporter,
	accountID int, audience *route.Domain, existingToken *models.RefreshToken,
) (string, string, error) {
	var err error
	err = SessionEnder(refreshTokenStore, existingToken)
	if err != nil {
		reporter.ReportError(errors.Wrap(err, "SessionEnder"))
	}

	// track actives
	if actives != nil {
		err = actives.Track(accountID)
		if err != nil {
			reporter.ReportError(errors.Wrap(err, "Track"))
		}
	}

	// track last activity
	_, err = accountStore.SetLastLogin(accountID)
	if err != nil {
		reporter.ReportError(errors.Wrap(err, "SetLastLogin"))
	}

	// create new session token
	session, err := sessions.New(refreshTokenStore, cfg, accountID, audience.String())
	if err != nil {
		return "", "", errors.Wrap(err, "sessions.New")
	}
	sessionToken, err := session.Sign(cfg.SessionSigningKey)
	if err != nil {
		return "", "", errors.Wrap(err, "session.Sign")
	}

	// create new identity token
	identityToken, err := identities.New(cfg, session, accountID, audience.String()).Sign(keyStore.Key())
	if err != nil {
		return "", "", errors.Wrap(err, "identities.New")
	}

	return sessionToken, identityToken, nil
}

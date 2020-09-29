package commands

import (
	"fmt"

	"code.cloudfoundry.org/credhub-cli/config"
	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth/uaa"
	"code.cloudfoundry.org/credhub-cli/util"
)

type LogoutCommand struct {
	ConfigCommand
}

func (c *LogoutCommand) Execute([]string) error {
	if err := RevokeTokenIfNecessary(c.config); err != nil {
		return err
	}
	MarkTokensAsRevokedInConfig(&c.config)
	if err := config.WriteConfig(c.config); err != nil {
		return err
	}
	fmt.Println("Logout Successful")
	return nil
}

func RevokeTokenIfNecessary(cfg config.Config) error {
	credhubClient, err := credhub.New(cfg.ApiURL, credhub.CaCerts(cfg.CaCerts...), credhub.SkipTLSValidation(cfg.InsecureSkipVerify), credhub.SetHttpTimeout(cfg.HttpTimeout))
	if err != nil {
		return err
	}

	uaaClient := uaa.Client{
		AuthURL: cfg.AuthURL,
		Client:  credhubClient.Client(),
	}

	if util.TokenIsPresent(cfg.AccessToken) {
		return uaaClient.RevokeToken(cfg.AccessToken)
	}

	return nil
}

func MarkTokensAsRevokedInConfig(cfg *config.Config) {
	cfg.AccessToken = "revoked"
	cfg.RefreshToken = "revoked"
}

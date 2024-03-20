package commands

import (
	"encoding/json"

	"code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/errors"
)

type RegenerateCommand struct {
	CredentialIdentifier string `required:"yes" short:"n" long:"name" description:"Selects the credential to regenerate"`
	Metadata             string `long:"metadata" description:"[JSON] Sets additional metadata on the credential"`
	OutputJSON           bool   `short:"j" long:"output-json" description:"Return response in JSON format"`
	ClientCommand
}

func (c *RegenerateCommand) Execute([]string) error {
	var options []credhub.RegenerateOption
	if c.Metadata != "" {
		var metadata credentials.Metadata
		if err := json.Unmarshal([]byte(c.Metadata), &metadata); err != nil {
			return errors.NewInvalidJSONMetadataError()
		}

		options = appendMetadataOptions(metadata, options)
	} else {
		credential, err := c.client.GetLatestVersion(c.CredentialIdentifier)
		if err != nil {
			return err
		}

		options = appendMetadataOptions(credential.Metadata, options)
	}

	credential, err := c.client.Regenerate(c.CredentialIdentifier, options...)

	if err == credhub.ServerDoesNotSupportMetadataError {
		return errors.NewServerDoesNotSupportMetadataError()
	}

	if err != nil {
		return err
	}

	credential.Value = "<redacted>"
	formatOutput(c.OutputJSON, credential)

	return nil
}

func appendMetadataOptions(metadata credentials.Metadata, options []credhub.RegenerateOption) []credhub.RegenerateOption {
	withMetadata := func(g *credhub.RegenerateOptions) error {
		g.Metadata = metadata
		return nil
	}

	options = append(options, withMetadata)
	return options
}

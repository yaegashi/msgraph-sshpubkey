package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/yaegashi/msgraph.go/msauth"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
	"golang.org/x/oauth2"
)

// KeyCmd is get/set/delete subcommand
type KeyCmd struct {
	*Cmd
	In          string                              `json:"-"`
	Out         string                              `json:"-"`
	Login       string                              `json:"login,omitempty"`
	TokenSource oauth2.TokenSource                  `json:"-"`
	GraphClient *msgraph.GraphServiceRequestBuilder `json:"-"`
}

// User returns target graph user endpoint
func (cmd *KeyCmd) User() *msgraph.UserRequestBuilder {
	if cmd.Login == "" {
		return cmd.GraphClient.Me()
	}
	userID := cmd.Login
	if cmd.LoginMap != nil {
		if id, ok := cmd.LoginMap[cmd.Login]; ok {
			userID = id
		}
	}
	return cmd.GraphClient.Users().ID(userID)
}

// Authenticate performs OAuth2 authentication
func (cmd *KeyCmd) Authenticate(ctx context.Context) error {
	var err error
	m := msauth.NewManager()
	if cmd.ClientSecret == "" {
		cmd.TokenSource, err = m.DeviceAuthorizationGrant(ctx, cmd.TenantID, cmd.ClientID, defaultScopes, nil)
		if err != nil {
			return err
		}
	} else {
		scopes := []string{msauth.DefaultMSGraphScope}
		cmd.TokenSource, err = m.ClientCredentialsGrant(ctx, cmd.TenantID, cmd.ClientID, cmd.ClientSecret, scopes)
		if err != nil {
			return err
		}
	}
	cmd.GraphClient = msgraph.NewClient(oauth2.NewClient(ctx, cmd.TokenSource))
	return nil
}

// Get performs get operation on user's extensions
func (cmd *KeyCmd) Get(args []string) error {
	ctx := context.Background()
	err := cmd.Authenticate(ctx)
	if err != nil {
		return err
	}
	r := cmd.User().Request()
	r.Select("id")
	r.Expand("extensions")
	user, err := r.Get(ctx)
	if err != nil {
		return err
	}
	for _, x := range user.Extensions {
		if *x.ID != cmd.ExtensionName {
			continue
		}
		value, _ := x.GetAdditionalData("value")
		if s, ok := value.(string); ok {
			if cmd.Out == "-" {
				fmt.Print(s)
				return nil
			}
			return ioutil.WriteFile(cmd.Out, []byte(s), 0644)
		}
		return fmt.Errorf("No value in extension %s", cmd.ExtensionName)
	}
	return nil
}

// Set performs set operation on user's extension
func (cmd *KeyCmd) Set(args []string) error {
	ctx := context.Background()
	err := cmd.Authenticate(ctx)
	if err != nil {
		return err
	}
	var in []byte
	if cmd.In == "-" {
		in, err = ioutil.ReadAll(os.Stdin)
	} else {
		in, err = ioutil.ReadFile(cmd.In)
	}
	if err != nil {
		return err
	}
	newExt := &msgraph.Extension{}
	newExt.SetAdditionalData("extensionName", cmd.ExtensionName)
	newExt.SetAdditionalData("value", string(in))
	_, err = cmd.User().Extensions().Request().Add(ctx, newExt)
	if err != nil {
		if errRes, ok := err.(*msgraph.ErrorResponse); ok {
			if errRes.StatusCode() == http.StatusConflict {
				err = cmd.User().Extensions().ID(cmd.ExtensionName).Request().Update(ctx, newExt)
			}
		}
	}
	return err
}

// Delete performs delete operation on user's extension
func (cmd *KeyCmd) Delete(args []string) error {
	ctx := context.Background()
	err := cmd.Authenticate(ctx)
	if err != nil {
		return err
	}
	return cmd.User().Extensions().ID(cmd.ExtensionName).Request().Delete(ctx)
}

// Main is KeyCmd entrypoint
func (cmd *KeyCmd) Main(args []string) error {
	cmd.FlagSet.StringVar(&cmd.Login, "login", "", "login name")
	cmd.FlagSet.StringVar(&cmd.In, "in", "", "Input file (\"-\" for stdin)")
	cmd.FlagSet.StringVar(&cmd.Out, "out", "-", "Output file (\"-\" for stdout)")
	cmd.FlagSet.Parse(args[1:])
	err := cmd.LoadConfig(&cmd)
	if err != nil {
		return err
	}
	arg0 := args[0]
	args = cmd.FlagSet.Args()
	switch arg0 {
	case "get":
		return cmd.Get(args)
	case "set":
		return cmd.Set(args)
	case "delete":
		return cmd.Delete(args)
	}
	return fmt.Errorf("Unknown command: %s", arg0)
}

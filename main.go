package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	defaultExtensionName = "dev.l0w.ssh_public_keys"
	defaultTenantID      = "common"
	defaultClientID      = "45c7f99c-0a94-42ff-a6d8-a8d657229e8c"
)

// Cmd is main command
type Cmd struct {
	FlagSet       *flag.FlagSet     `json:"-"`
	Config        string            `json:"-"`
	TenantID      string            `json:"tenant_id,omitempty"`
	ClientID      string            `json:"client_id,omitempty"`
	ClientSecret  string            `json:"client_secret,omitempty"`
	ExtensionName string            `json:"extension_name,omitempty"`
	LoginMap      map[string]string `json:"login_map,omitempty"`
}

// Main is MainCmd entrypoint
func (cmd *Cmd) Main(args []string) error {
	cmd.FlagSet = flag.NewFlagSet(args[0], flag.ExitOnError)
	cmd.FlagSet.Usage = func() {
		out := cmd.FlagSet.Output()
		fmt.Fprintf(out, "Usage: %s <command> [options]\n", cmd.FlagSet.Name())
		fmt.Fprintf(out, "Commands:\n")
		fmt.Fprintf(out, "  get    - Get SSH public keys\n")
		fmt.Fprintf(out, "  set    - Set SSH public keys\n")
		fmt.Fprintf(out, "  delete - Delete SSH public keys\n")
		fmt.Fprintf(out, "  serve  - Start web server\n")
		fmt.Fprintf(out, "Options:\n")
		cmd.FlagSet.PrintDefaults()
	}
	cmd.FlagSet.StringVar(&cmd.Config, "config", "", "Load config from a file")
	cmd.FlagSet.StringVar(&cmd.TenantID, "tenant-id", defaultTenantID, "Tenant ID")
	cmd.FlagSet.StringVar(&cmd.ClientID, "client-id", defaultClientID, "Client ID")
	cmd.FlagSet.StringVar(&cmd.ClientSecret, "client-secret", "", "Client secret (for client credentials grant)")
	cmd.FlagSet.StringVar(&cmd.ExtensionName, "extension-name", defaultExtensionName, "Extension name")
	cmd.FlagSet.Parse(args[1:])
	args = cmd.FlagSet.Args()
	if len(args) == 0 {
		return fmt.Errorf("Available command: get, set, serve")
	}
	switch args[0] {
	case "serve":
		serveCmd := &ServeCmd{Cmd: cmd}
		return serveCmd.Main(args)
	case "get", "set":
		keyCmd := &KeyCmd{Cmd: cmd}
		return keyCmd.Main(args)
	}
	return fmt.Errorf("Unknown command: %s", args[0])
}

// LoadConfig loads config from a file
func (cmd *Cmd) LoadConfig(obj interface{}) error {
	if cmd.Config == "" {
		return nil
	}
	b, err := ioutil.ReadFile(cmd.Config)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

func main() {
	app := &Cmd{}
	err := app.Main(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var defaultScopes = []string{"openid", "offline_access", "User.ReadWrite"}

type HTTPServerError struct {
	Code int
	Text string
}

func (err *HTTPServerError) Error() string {
	return fmt.Sprintf("%d %s\n%s\n", err.Code, http.StatusText(err.Code), err.Text)
}

type spaHandler string

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	root := string(h)
	path := filepath.Join(root, filepath.Clean(r.URL.Path))
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		http.ServeFile(w, r, filepath.Join(root, "index.html"))
		return
	}
	http.ServeFile(w, r, path)
}

// ServeCmd is serve subcommand
type ServeCmd struct {
	*Cmd
	Listen       string         `json:"listen,omitempty"`
	Root         string         `json:"root,omitempty"`
	RedirectURI  string         `json:"redirect_uri,omitempty"`
	SessionDir   string         `json:"session_dir,omitempty"`
	SessionKey   string         `json:"session_key,omitempty"`
	SessionStore sessions.Store `json:"-"`
	OAuth2Config *oauth2.Config `json:"-"`
}

// Main is ServeCmd entrypoint
func (cmd *ServeCmd) Main(args []string) error {
	cmd.FlagSet.StringVar(&cmd.Listen, "listen", ":8080", "address:port to listen")
	cmd.FlagSet.StringVar(&cmd.Root, "root", "./app/build", "app root directory")
	cmd.FlagSet.StringVar(&cmd.RedirectURI, "redirect-uri", "http://localhost:8080/auth/callback", "Redirect URI")
	cmd.FlagSet.StringVar(&cmd.SessionDir, "session-dir", "", "session dir")
	cmd.FlagSet.StringVar(&cmd.SessionKey, "session-key", "", "session key")
	cmd.FlagSet.Parse(args[1:])
	err := cmd.LoadConfig(&cmd)
	if err != nil {
		return err
	}
	cmd.OAuth2Config = &oauth2.Config{
		ClientID:     cmd.ClientID,
		ClientSecret: cmd.ClientSecret,
		Endpoint:     microsoft.AzureADEndpoint(cmd.TenantID),
		RedirectURL:  cmd.RedirectURI,
		Scopes:       defaultScopes,
	}
	s := sessions.NewFilesystemStore(cmd.SessionDir, []byte(cmd.SessionKey))
	s.MaxLength(16384)
	cmd.SessionStore = s
	http.HandleFunc("/api/", cmd.HandleAPI)
	http.HandleFunc("/auth/", cmd.HandleAuth)
	http.Handle("/", spaHandler(cmd.Root))
	log.Printf("Server starting on %s", cmd.Listen)
	return http.ListenAndServe(cmd.Listen, nil)
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/okta/okta-sdk-golang/v5/okta"

	"golang.org/x/oauth2"
)

var (
	oktaDomain   = "https://{yourOktaDomain}"       // Replace with your Okta domain
	clientID     = "{yourClientID}"                 // Replace with your Client ID
	clientSecret = "{yourClientSecret}"             // Replace with your Client Secret
	redirectURI  = "http://localhost:8080/callback" // Replace with your redirect URI
	oauth2Config *oauth2.Config
)

func init() {
	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/oauth2/default/v1/authorize", oktaDomain),
			TokenURL: fmt.Sprintf("%s/oauth2/default/v1/token", oktaDomain),
		},
	}
}

func loginroutes() {
	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("GET")
	r.HandleFunc("/callback", callbackHandler).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code found", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	config, err := okta.NewConfiguration(
		okta.WithOrgUrl("https://{yourOktaDomain}"),
		okta.WithToken(token.AccessToken),
		okta.WithRequestTimeout(45),
		okta.WithRateLimitMaxRetries(3),
	)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	client := okta.NewAPIClient(config)

	user, resp, err := client.UserAPI.GetUser(client.GetConfig().Context, "{UserId|Username|Email}").Execute()
	if err != nil {
		fmt.Printf("Error Getting User: %v\n", err)
	}
	fmt.Printf("User: %+v\n Response: %+v\n\n", user, resp)

	fmt.Fprintf(w, "Hello, %s! Your email is %s.", user.Profile.FirstName, user.Profile.Email)
}

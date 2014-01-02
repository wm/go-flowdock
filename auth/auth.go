// Basic OAuth2 command line helper.
// This code is based on the example found at https://code.google.com/p/goauth2/source/browse/oauth/example/oauthreq.go
// See LICENSE-GOOGLE

package auth

import (
	"flag"
	"os"
	"log"
	"fmt"
	"net/http"
	"code.google.com/p/goauth2/oauth"
)

var (
	clientId     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client Secret")
	scope        = flag.String("scope", "flow private manage profile offline_access", "OAuth scope")
	redirectURL  = flag.String("redirect_url", "urn:ietf:wg:oauth:2.0:oob", "Redirect URL")
	authURL      = flag.String("auth_url", "https://api.flowdock.com/oauth/authorize", "Authentication URL")
	tokenURL     = flag.String("token_url", "https://api.flowdock.com/oauth/token", "Token URL")
	code         = flag.String("code", "", "Authorization Code")
	cachefile    = flag.String("cache", "cache.json", "Token cache file")
)

const usageMsg = `
To obtain a request token you must specify both -id and -secret.

To obtain Client ID and Secret, see the "OAuth 2 Credentials" section under
the "API Access" tab on this page: https://flowdock.com/account/authorized_applications

Once you have completed the OAuth flow, the credentials should be stored inside
the file specified by -cache and you may run without the -id and -secret flags.
`

func AuthenticationRequest() *http.Client {
	flag.Parse()

	// Set up a configuration.
	config := &oauth.Config{
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
		RedirectURL:  *redirectURL,
		Scope:        *scope,
		AuthURL:      *authURL,
		TokenURL:     *tokenURL,
		TokenCache:   oauth.CacheFile(*cachefile),
	}

	// Set up a Transport using the config.
	transport := &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if *clientId == "" || *clientSecret == "" {
			flag.Usage()
			fmt.Fprint(os.Stderr, usageMsg)
			os.Exit(2)
		}
		if *code == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("")
			fmt.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE\n")
			fmt.Println(url)
			os.Exit(0)
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(*code)
		if err != nil {
			log.Fatal("Exchange:", err)
		}
		// (The Exchange method will automatically cache the token.)
		fmt.Printf("Token is cached in %v\n", config.TokenCache)
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transport.Token = token
	client := transport.Client()

    return client
}

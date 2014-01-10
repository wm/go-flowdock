package main

import (
	"github.com/wm/go-flowdock/flowdock"
	"github.com/codegangsta/cli"
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"net/http"
	"log"
	"os"
	"sort"
)

const usageMsg = `
To obtain a request token you must specify both -id and -secret.

To obtain Client ID and Secret, see the "OAuth 2 Credentials" section under
the "API Access" tab on this page: https://flowdock.com/account/authorized_applications

Once you have completed the OAuth flow, the credentials should be stored inside
the file specified by -cache and you may run without the -id and -secret flags.
`

type Query struct {
     Tags   []string
     Org    string
     Flow   string
}

type AppDeployCount struct {
	Project       string
	Total         int
	DeployCount   *map[string]int
}

const limit = 100

// example of counting number of deploys per month for specidied applications.
//
// This is utilizing the fact that a deploy command messages the flow inbox and is tagged appropriately.
func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag {
		cli.StringFlag{"id", "", "Client ID"},
		cli.StringFlag{"secret", "", "Client Secret"},
		cli.StringFlag{"redirect_url", "urn:ietf:wg:oauth:2.0:oob", "Redirect URL"},
		cli.StringFlag{"auth_url", "https://api.flowdock.com/oauth/authorize", "Authentication URL"},
		cli.StringFlag{"token_url", "https://api.flowdock.com/oauth/token", "Token URL"},
		cli.StringFlag{"code", "", "Authorization Code"},
		cli.StringFlag{"cache", "cache.json", "Token cache file"},
		cli.StringFlag{"environment, e", "production", "the deploy target"},
		cli.StringFlag{"organization, o", "iora", "the organization of the flow"},
		cli.StringFlag{"flow, f", "tech-stuff", "the name of the flow to query"},
	}

	app.Name = "deploys"
	app.Usage = "Counts the deploys for the listed applications by month"

	app.Action = func(c *cli.Context) {
		client := flowdock.NewClient(AuthenticationRequest(c))
		// args := []string{"bouncah", "icis", "cronos", "snowflake"} //c.Args()
		args := c.Args()

		if len(args) == 0 {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		channel := make(chan AppDeployCount, len([]string(args)))

		for _, app := range []string(args) {
			tags   := []string{
				"deployment",
				"deploy_end",
				c.String("environment"),
				app,
			}
			q := Query{tags, c.String("organization"), c.String("flow")}
			getAppDeployCount(q, client, channel)
		}

		displayAppDeployCount(channel)
	}

	app.Run(os.Args)
}

func getAppDeployCount(q Query, client *flowdock.Client, channel chan AppDeployCount) {
	var deployCount = map[string]int{}

	go func() {

		app := q.Tags[len(q.Tags)-1]
		opt := flowdock.MessagesListOptions{Limit: 100, TagMode: "and"}

		opt.Tags   = q.Tags
		opt.Search = "production to production"
		opt.Event  = "mail"

		messages, _, err := client.Messages.List(q.Org, q.Flow, &opt)

		if err != nil {
			log.Fatal("Get:", err)
		}

		total := 0
		for _, msg := range messages {
			if !stringInSlice("preproduction", *msg.Tags) {
				total++
				month := msg.Sent.Format("2006-Jan")
				deployCount[month]++
				// fmt.Println("MSG:", month, *msg.ID, *msg.Event, *msg.Tags)
			}
		}

		if len(messages) == limit {
			removeEarliestMonth(&deployCount)
		}

		channel <- AppDeployCount{app, total, &deployCount} 
	}()
}

func removeEarliestMonth(displayCount *map[string]int) {
	sortedKeys := sortedKeys(displayCount)
	firstKey   := (*sortedKeys)[0]
	delete(*displayCount, firstKey)
}

func displayAppDeployCount(adcChan <-chan AppDeployCount) {
	for i := 0; i < cap(adcChan); i++ {
		adc := <-adcChan

		fmt.Println()
		fmt.Println("Application:", adc.Project)
		fmt.Println()

		for k, v := range *adc.DeployCount {
			fmt.Println(k, v)
		}

		fmt.Println()
		fmt.Println("  Total:", adc.Total)
		fmt.Println()
	}
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func sortedKeys(m *map[string]int) *[]string  {
    mk := make([]string, len(*m))
    i := 0
    for k, _ := range *m {
        mk[i] = k
        i++
    }
	sort.Strings(mk)
	return &mk
}

func AuthenticationRequest(c *cli.Context) *http.Client {
	// Set up a configuration.
	config := &oauth.Config{
		ClientId:     c.String("id"),
		ClientSecret: c.String("secret"),
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scope:        "",
		AuthURL:      c.String("auth_url"),
		TokenURL:     c.String("token_url"),
		TokenCache:   oauth.CacheFile(c.String("cache")),
	}

	// Set up a Transport using the config.
	transport := &oauth.Transport{Config: config}

	// Try to pull the token from the cache; if this fails, we need to get one.
	token, err := config.TokenCache.Token()
	if err != nil {
		if c.String("id") == "" || c.String("secret") == "" {
			cli.ShowAppHelp(c)
			fmt.Fprint(os.Stderr, usageMsg)
			os.Exit(2)
		}
		if c.String("code") == "" {
			// Get an authorization code from the data provider.
			// ("Please ask the user if I can access this resource.")
			url := config.AuthCodeURL("")
			fmt.Println("Visit this URL to get a code, then run again with -code=YOUR_CODE\n")
			fmt.Println(url)
			os.Exit(0)
		}
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		token, err = transport.Exchange(c.String("code"))
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

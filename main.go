package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const name = "gh-star"

const version = "0.0.1"

var revision = "HEAD"

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "V", false, "Print the version")
	if showVersion {
		fmt.Printf("%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return
	}

	b, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		os.Exit(1)
	}
	urls := map[string]struct{}{}
	for _, line := range strings.Split(string(b), "\n") {
		tok := strings.Fields(line)
		if len(tok) < 2 {
			continue
		}
		urls[tok[1]] = struct{}{}
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	for k, _ := range urls {
		u, err := url.Parse(k)
		if err != nil {
			continue
		}
		p := strings.Trim(u.Path, "/")
		p = strings.TrimSuffix(p, ".git")
		tok := strings.Split(p, "/")
		_, err = client.Activity.Star(context.Background(), tok[0], tok[1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("★ %s/%s\n", tok[0], tok[1])
	}
}

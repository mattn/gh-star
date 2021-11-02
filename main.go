package main

import (
	"bytes"
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
	flag.Parse()

	if showVersion {
		fmt.Printf("%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return
	}

	token := os.Getenv("GITHUB_TOKEN")
	if len(token) == 0 {
		fmt.Fprintln(os.Stderr, "GITHUB_TOKEN is not set")
		os.Exit(1)
	}

	urls := map[string]struct{}{}

	if flag.NArg() == 0 {
		var buf bytes.Buffer
		cmd := exec.Command("git", "remote", "-v")
		cmd.Stderr = &buf
		b, err := cmd.Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, buf.String())
			os.Exit(1)
		}
		for _, line := range strings.Split(string(b), "\n") {
			tok := strings.Fields(line)
			if len(tok) < 2 {
				continue
			}
			urls[tok[1]] = struct{}{}
		}
	} else {
		for _, arg := range flag.Args() {
			urls[arg] = struct{}{}
		}
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
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

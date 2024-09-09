/*
This app shares a folder via HTTP.
Useful for quick sharing. Not suitable for public hosting over the internet.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

const (
	appName    = "sharefolder"
	appVersion = "v0.2.1"
	appAuthor  = "Andras Belicza/Jm61"
	appHome    = "https://github.com/icza/toolbox"
)

var (
	version        = flag.Bool("version", false, "print version info and exit")
	addr           = flag.String("addr", ":8080", "address to start the server on")
	username	   = flag.String("username", "", "require basic authentication username")
	password       = flag.String("password", "", "require basic authentication password")
	promptPassword = flag.Bool("promptPassword", false, "prompt for password to enter in the console if you don't want to provide with -password")
)

func main() {
	flag.Usage = printUsage
	flag.Parse()

	if *version {
		printVersion()
		return
	}

	if *promptPassword {
		fmt.Print("Enter basic authentication password: ")
		scanner := bufio.NewScanner(os.Stdout)
		scanner.Scan()
		*password = scanner.Text()
	}

	if *password != "" {
		log.Print("Using basic auth password ", strings.Repeat("*", utf8.RuneCountInString(*password)))
	}
	if *username != "" {
		log.Print("Using basic auth username ", strings.Repeat("*", utf8.RuneCountInString(*username)))
	}

	args := flag.Args()

	path := ""
	if len(args) > 0 {
		path = args[0]
	}

	path, err := filepath.Abs(path)
	if err != nil {
		log.Printf("Failed to resolve %s", path)
		os.Exit(1)
	}

	log.Printf("Sharing folder: %s", path)

	// Find out and print which addresses we're listening on:
	host, port, err := net.SplitHostPort(*addr)
	if err != nil {
		log.Printf("Failed to split addr: %v", err)
		os.Exit(2)
	}
	if host != "" {
		// Host is explicit:
		log.Printf("Listening on http://%s/", *addr)
	} else {
		// Host is missing, we'll listen on all available interfaces:
		printLocalInterfaces(port)
	}

	fsys := dotFileHidingFileSystem{http.Dir(path)}
	http.Handle("/", basicAuth(http.FileServer(fsys)))
	log.Print(http.ListenAndServe(*addr, nil))
}



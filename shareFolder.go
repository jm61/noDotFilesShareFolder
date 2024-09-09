package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// dotFileHidingFileSystem is an http.FileSystem that hides
// hidden "dot files" from being served.
type dotFileHidingFileSystem struct {
	http.FileSystem
}
// dotFileHidingFile is the http.File use in dotFileHidingFileSystem.
// It is used to wrap the Readdir method of http.File so that we can
// remove files and directories that start with a period from its output.
type dotFileHidingFile struct {
	http.File
}

func printVersion() {
	fmt.Println(appName, "version:", appVersion)
	fmt.Println("Platform:", runtime.GOOS, runtime.GOARCH)
	fmt.Println("Built with:", runtime.Version())
	fmt.Println("Author:", appAuthor)
	fmt.Println("Home page:", appHome)
}

func printUsage() {
	fmt.Println("Usage:")
	name := os.Args[0]
	fmt.Printf("%s [FLAGS] [folder-to-share]\n", name)
	fmt.Println("(The current working directory is shared if not specified.)")
	fmt.Println()
	fmt.Println("Flags:")

	flag.PrintDefaults()
}

func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, pass, ok := r.BasicAuth(); !ok || // Missing / invalid basic auth
			pass != *password || user != *username { // Invalid password
			w.Header().Set("WWW-Authenticate", `Basic realm="sharefolder"`)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func printLocalInterfaces(port string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Failed to get interfaces: %v", err)
		os.Exit(11)
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Printf("Failed to get interface addresses: %v", err)
			os.Exit(12)
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ipv4 := ipNet.IP.To4(); ipv4 != nil {
					log.Printf("Listening on http://%s:%s/", ipv4, port)
				}
			}
		}
	}
}

// containsDotFile reports whether name contains a path element starting with a period.
// The name is assumed to be a delimited by forward slashes, as guaranteed
// by the http.FileSystem interface.
func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}
	return false
}

// Readdir is a wrapper around the Readdir method of the embedded File
// that filters out all files that start with a period in their name.
func (f dotFileHidingFile) Readdir(n int) (fis []fs.FileInfo, err error) {
	files, err := f.File.Readdir(n)
	for _, file := range files { // Filters out the dot files
		if !strings.HasPrefix(file.Name(), ".") {
			fis = append(fis, file)
		}
	}
	return
}

// Open is a wrapper around the Open method of the embedded FileSystem
// that serves a 403 permission error when name has a file or directory
// with whose name starts with a period in its path.
func (fsys dotFileHidingFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) { // If dot file, return 403 response
		return nil, fs.ErrPermission
	}

	file, err := fsys.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return dotFileHidingFile{file}, err
}
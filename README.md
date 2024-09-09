# sharefolder

This app shares a folder via HTTP.

Useful for quick sharing. Not suitable for public hosting over the internet.

## Usage

    sharefolder [FLAGS] [folder-to-share]

The current working directory is shared if not specified.

### Initial repo

https://github.com/icza/toolbox/tree/main/cmd/sharefolder

### Added noDotFiles part

https://pkg.go.dev/net/http#example-FileServer-DotFileHiding

### Note about http Dir type

Note that Dir could expose sensitive files and directories. Dir will follow symlinks pointing out of the directory tree, which can be especially dangerous if serving from a directory in which users are able to create arbitrary symlinks. Dir will also allow access to files and directories starting with a period, which could expose sensitive directories like .git or sensitive files like .htpasswd. To exclude files with a leading period, remove the files/directories from the server or create a custom FileSystem implementation.

# linvestor

## A finance application for the share club Linvestor.

### Installing GO:
- Download [GO](http://golang.org/dl/)
 - Install the .MSI package WITH ALL DEFAULT. Just hit NEXT.
- Create folder "C:\GoPath"
- Go to the "System" control panel, click the "Advanced" tab. Select "Environment Variables" and under "System variables":
 - Add GOPATH variable, set it to "C:\GoPath"

### Install go tools
- Open a NEW terminal (with new env vars)
- Install good tools:
 - `$ go get -u -ldflags -H=windowsgui github.com/nsf/gocode/...`
 - `$ go get -u github.com/abourget/godef/...`
 - `$ go get -u golang.org/x/tools/cmd/...`

### Install Atom (IDE)
- Download [Atom](https://atom.io/)
- Install and then open
- In the Welcome Guide
 - Click "Install package"
 - Click "Open Installer"
 - Search "go-plus" (by joefitzgerald)
  - Install it

### Install Git (version control)
- Download [Git](http://git-scm.com/downloads)
- Install
 - When asked to select how to "Adjust your PATH environment", select "Use Git from the Windows Command Prompt".
 - Continue on with the defaults.

### Download the code to local environment
- `$ go get github.com/fabwi987/linvestor`
 - The project will download to C:\GoPath\src\github.com\fabwi987\linvestor







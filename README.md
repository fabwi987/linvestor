# linvestor

## A finance application for the share club Linvestor made in GO.

## Getting started guide
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

### Using Heoku (hosting through PAAS)
- Create a [Heroku account](https://heroku.com)
- Download and install [Heroku tool belt}(https://toolbelt.heroku.com/)
- Start a cmd-window from the project root
 - `$ heroku login`
 - `$ godep save ./..`
 - `$ go install ./..`
 - `$ heroku local`
- Open a browser and go to "localhost:[port]/" to run the application locally

### Help
- [Learning GO](https://golang.org/doc/)
- [Using GIT](http://rogerdudler.github.io/git-guide/)
- [Heroku and GO](https://devcenter.heroku.com/articles/getting-started-with-go#introduction)
- [Atom tips](http://readwrite.com/2014/05/20/github-atom-5-tips-getting-started-tutorial-corey-johnson/)







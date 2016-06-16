# linvestor

Getting started:

Go to:
http://golang.org/dl/
and download and install the .MSI package WITH ALL DEFAULT. Just hit NEXT.
Create new folder in C:\ named "GoPath"
Go to the "System" control panel, click the "Advanced" tab. Select "Environment Variables" and under "System variables":
add GOPATH variable, set it to "C:\GoPath"

Install Git (version control)
Go to:
http://git-scm.com/downloads
Download and install, when asked to select how to "Adjust your PATH environment", select the SECOND choice which is "Use Git from the Windows Command Prompt".
Continue on with the defaults.

Install go tools
Open a NEW terminal (with new env vars), paste this in to install good tools:
go get -u -ldflags -H=windowsgui github.com/nsf/gocode/...
go get -u github.com/abourget/godef/...
go get -u golang.org/x/tools/cmd/...
You might see an error with godoc. Ignore it.

Install Atom (IDE)
Atom has the best support for coding in Go, has excellent plugins and works across platforms. Throw Sublime Text away and use Atom.
Go to:
https://atom.io/
Download and install (it will open), and once in there:
In the Welcome Guide
Click "Install package"
Click "Open Installer"
Search "go-plus" (by joefitzgerald)
Install it

Getting the repository from git:
write in cmd: go get github.com/fabwi987/linvestor
project will download to C:\GoPath\src\github.com\fabwi987\linvestor

Installing Heroku:
Download Heroku toolbelt: https://devcenter.heroku.com/articles/getting-started-with-go#set-up
From cmd: heroku login (Get user account and password from admin)
Associate heroku with the app: $heroku git:remote -a [name_of_application]


Testing (all commands are run from the root directory of "linvestor"):
godep save ./..
go install ./..
heroku local
Open browser -> loclahost:5050

GoLive
Add the modified files to the local git repository: $ git add -A .
Commit the changes to the repository: $ git commit -m "[Insert message]"
Deploy: $ git push heroku master







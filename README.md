# What is Valorant Account Tracker?

This app allows you to track Valorant accounts you have used locally. The main
purpose is to store the usernames and last known competitive rank in order to
make picking an account easier.\
For example, if you have played on multiple accounts and do not remember what
rank each of your accounts ended up in, you can use this app to check.

## Roadmap:

<pre>
[x] Simple UI
  [ ] Notes functionality
  [ ] Automatic refreshes
[x] Buildable executable with no dependencies
  [ ] Automatically test and build on push
[x] A persistent service running in the background
[x] Automatically detect when the account is changed and pull rank details
</pre>

## How to get started

- Download the
  [latest release](https://github.com/Cufee/valorant-account-tracker-go/releases)
  \
  **OR**
- Clone the code
- Install [go](https://go.dev/) and [Task](https://taskfile.dev/installation/)
- `task build`
- Open Valorant
- Run the executable in `build/`
- (Optional) Add to Windows startup folder

## How does it work?

Riot Client opens up a web server each time the app is launched. This API
exposes some basic information about the running instance of Valorant and the
account you are currently using.\
Using this web server, we are able to request tokens for Riot API, using which
we can collect all sorts of information. As far as I understand, we can access
the entirety of the [Valorant API](https://valapidocs.techchrism.me/).

Rewrite of my [draft project](https://github.com/Cufee/valorant-account-tracker)
made with Deno

**This app is a side project of mine and is in very early stages. I am planning
on further refining this concept into something that can be easily used by
anyone.**

Here is what I have in mind:\
[x] A persistent service running in the background\
[ ] Automatically detect when the account is changed and pull rank details\
[x] Simple UI

# How to get started

- Clone the code
- Install [go](https://go.dev/)
- `go build -o build/valorantAccountTracker.exe main.go`
- Open Valorant
- Run the executable in `build/`
- (Optional) Add to Windows startup folder

# What is Valorant Account Tracker?

This app allows you to track Valorant accounts you have used locally. The main
purpose is to store the usernames and last known competitive rank in order to
make picking an account easier.

For example, if you have played on multiple accounts and do not remember what
rank each of your accounts ended up in, you can use this app to check.

## How does it work?

Riot Client opens up a web server each time the app is launched. This API
exposes some basic information about the running instance of Valorant and the
account you are currently using.

Using this web server, we are able to request tokens for Riot API, using which
we can collect all sorts of information. As far as I understand, we can access
the entirety of the [Valorant API](https://valapidocs.techchrism.me/).

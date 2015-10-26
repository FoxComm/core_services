# FoxComm Core Services

## Dependencies

1. Go 1.5.1: `brew install go`

2. Vulcand: `go get github.com/FoxComm/vulcand`

## Setting up your environment

1. Get GPM: `brew install gpm`; instructions for [Non-OSX setups here.](https://github.com/pote/gpm)

2. Get dependencies via GPM: `gpm install`  

## Common Issues

1. `go get github.com/FoxComm/vulcand` produces an error such as `coud not read username for`

Both gpm and `go get` support using private GitHub repositories! Here's what you need to do in order for a specific machine to be able to access them:

* Generate a GitHub access token by following [these instructions](https://help.github.com/articles/creating-an-access-token-for-command-line-use/).
* Add the following line to the `~/.netrc` file in your home directory.

```bash
machine github.com login <token>
```

You can now use gpm (and `go get`) to install private repositories to which your user has access! :)



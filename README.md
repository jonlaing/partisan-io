# Partisan.IO

> Social networking app for discussion (arguing) regarding politics.

## Prerequisites

### Golang

Installation instructions: [https://golang.org/doc/install](https://golang.org/doc/install)

### Homebrew

Run the following in your Terminal:

```bash
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

### NodeJS

Run the following in your Terminal:

```bash
brew install node
```

## Running your project

The project has two components, the Go backend and the JS front end. The repo should have a build of the backend as an executable, so to start it, run the following in your Terminal:

```bash
./partisan
```

If you would like to build the backend from source, you'll first need the Golang environment installed. You'll also need to run `go get` to get all the dependencies. After you've done that, run the following in your Terminal:

```bash
go build
```

And then run the executable as in the previous command. The server will run at `http://localhost:4000`.

To build the front end, simply run the following.

```bash
npm start
```

This will watch the `/src` directory for any changes and automatically compile them into the `/dist` directory, provided there are no errors. Live reloading currently does not work as it conflicts with the backend, but we may fix that issue shortly.


## Generating Additional Code

You can add additional functionality to your application by invoking the subgenerators included in the Flux Generator. You can add components using the following commands:

#### Components
```bash
$ yo flux:component ComponentName
```

#### Actions
```bash
$ yo flux:action ActionCreatorName
```

#### Stores
```bash
$ yo flux:store StoreName
```

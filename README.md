# Tania API Server

[![Build
Status](https://travis-ci.com/Tanibox/tania-server.svg?token=rFvCvk3YAEKpm9H4axVv&branch=master)](https://travis-ci.com/Tanibox/tania-server)

The future of Tania farm management system is here. 

## Requirement
- Go v1.9 
- Vue 2.x

## Installation
- Make sure you have installed `golang/dep` 
    - https://golang.github.io/dep/docs/installation.html
    - https://gist.github.com/subfuzion/12342599e26f5094e4e2d08e9d4ad50d
- Clone the repo using `go get github.com/Tanibox/tania-server`
- From the project root, call `dep ensure` to install the Go dependencies
    - If you have an issue with `dep ensure`, you can call `go get` instead.
- Create a new file `conf.json` using the values from the `conf.json.example` and set it with your own values.
- Call `npm install` to install Vue dependencies
- Call `npm run dev` to build the Vue
- Setup SQLite:
    - Edit `SqlitePath` in `conf.json` to your sqlite DB file path (ex: /Users/user/Programs/sqlite/tania.db)
    - Create empty file with the exact filename and path that match the `SqlitePath` config.
- Run the Go server using `go run main.go` and open it in the `http://localhost:8080`

## Test
- Call `go test ./...` to run all the Go tests.
- Call `npm run cypress:run` to run the end-to-end test

## License

Tania is available under Apache 2.0 open source license.

<div align="center">
    <img src="logo.png" alt="Tania Farm Management System" width="200">
    <h1>Farm Management Software</h1>
    <img src="https://img.shields.io/badge/semver-1.5.1-green.svg?maxAge=2592000" alt="semver">
    <img src="https://travis-ci.com/Tanibox/tania-core.svg?branch=master" alt="Build Status">
    <a href="https://opensource.org/licenses/Apache-2.0" target="_blank"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License"></a>
</div>


**Tania** is a free and open source farm management software. You can manage your farm areas, farm reservoirs, farm tasks, inventories, and the crop growing progress. It is designed for any type of farms.

## Requirements
- Go 1.9 
- Vue 2.x
- Node 8.9.x

## Installation
- Make sure you have installed `golang/dep` 
    - https://golang.github.io/dep/docs/installation.html
    - https://gist.github.com/subfuzion/12342599e26f5094e4e2d08e9d4ad50d
- Clone the repo using `go get github.com/Tanibox/tania-core`
- From the project root, call `dep ensure` to install the Go dependencies
    - If you have an issue with `dep ensure`, you can call `go get` instead.
- Create a new file `conf.json` using the values from the `conf.json.example` and set it with your own values.
- Call `npm install` to install Vue dependencies
- Call `npm run dev` to build the Vue
- Setup SQLite:
    - Edit `SqlitePath` in `conf.json` to your sqlite DB file path (ex: /Users/user/Programs/sqlite/tania.db)
    - Create empty file with the exact filename and path that match the `SqlitePath` config.
- Run the Go server using `go run main.go` and open it in the `http://localhost:8080`
- Default username and password are `tania / tania`

## Database Engine

Tania use SQLite as the default database engine. You may use MySQL as your database engine by replacing `sqlite` with `mysql` at `tania_persistence_engine` field in your `conf.json`.

```
{
    "app_port": "8080",
    "tania_persistence_engine": "mysql",
    "demo_mode": true,
    "upload_path_area": "uploads/areas",
    "upload_path_crop": "uploads/crops",
    "sqlite_path": "db/sqlite/tania.db",
    "mysql_host": "localhost",
    "mysql_port": "3306",
    "mysql_dbname": "your_mysql_db_name",
    "mysql_user": "your_mysql_user",
    "mysql_password": "your_mysql_password",
    "redirect_uri": [
        "http://localhost:8080",
        "http://127.0.0.1:8080"
    ],
    "client_id": "f0ece679-3f53-463e-b624-73e83049d6ac"
}
```

## Test
- Call `go test ./...` to run all the Go tests.
- Call `npm run cypress:run` to run the end-to-end test

## REST API
**Tania** have REST APIs to easily integrate with any softwares, even you can build a mobile app client for it. You can read the documentation here: [Tania REST API](https://documenter.getpostman.com/view/3434975/tania/RVnb9H2z).

## License

Tania is available under Apache 2.0 open source license.

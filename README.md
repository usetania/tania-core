<div align="center">
    <img src="logo.png" alt="Tania The Farmer Journal" width="200">
    <h1>The Farmer Journal</h1>
    <img src="https://opencollective.com/tania/tiers/backer/badge.svg?label=backer&color=brightgreen" />
    <img src="https://opencollective.com/tania/tiers/sponsor/badge.svg?label=sponsor&color=brightgreen" />
    <img src="https://img.shields.io/badge/semver-2.0.0-green.svg?maxAge=2592000" alt="semver">
    <a href="https://opensource.org/licenses/Apache-2.0" target="_blank"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License"></a>
</div>

# Warning

This is the development branch of Tania. Changes can occur nightly. If you need the stable branch you can checkout [the master branch](https://github.com/Tanibox/tania-core/tree/master).

## Roadmap

You can check the roadmap in [Tania's GitHub project](https://github.com/orgs/usetania/projects/6/views/1).

---

**Tania** is a free and open source farm management software. You can manage your farm areas, farm reservoirs, farm tasks, inventories, and the crop growing progress. It is designed for any type of farms.

Download Tania for Windows x64 and Linux x64 on [the release page](https://github.com/Tanibox/tania-core/releases/tag/1.7.1).

![Screenshot](screenshot.PNG)

## Getting Started

This software is built with [Go](https://golang.org) programming language. It means you will get an executable binary to run on your machine. You **don't need** extra software like MAMP, XAMPP, or WAMP to run **Tania**, but you may need MySQL database if you choose to use it instead of SQLite *(the default database.)*

If your OS is not listed on our releases page, you have to build Tania for your OS by yourself. You can follow our instructions to build **Tania**.

### Prerequisites
- [Go](https://golang.org) >= 1.16
- [NodeJS](https://nodejs.org/en/) >= 16

### Building Instructions

**THIS DOCUMENTATION WILL BE UPDATED LATER**

We are in the progress of building the new frontend application.

### Database Engine

Tania uses SQLite as the default database engine. You may use MySQL as your database engine by replacing `sqlite` with `mysql` at `tania_persistence_engine` field in your `backend/conf.json`.

```
{
  "app_port": "8080",
  "tania_persistence_engine": "sqlite",
  "demo_mode": true,
  "upload_path_area": "uploads/areas",
  "upload_path_crop": "uploads/crops",
  "sqlite_path": "db/sqlite/tania.db",
  "mysql_host": "127.0.0.1",
  "mysql_port": "3306",
  "mysql_dbname": "tania",
  "mysql_user": "root",
  "mysql_password": "root",
  "redirect_uri": [
      "http://localhost:8080",
      "http://127.0.0.1:8080"
  ],
  "client_id": "f0ece679-3f53-463e-b624-73e83049d6ac"
}
```

### Run The Test

Use `go test ./...` inside the `backend` folder to run all the Go tests.

## REST APIs
**Tania** have REST APIs to easily integrate with any softwares, even you can build a mobile app client for it. You can import the JSON file inside Postman directory to [Postman app](https://www.getpostman.com).

## Contributing to Tania

We welcome contributions, but request you to follow these [guidelines](contributing.md).

### Localisation

You can help us to localise Tania into your language by following these steps:

1. Copy `frontend/languages/template.pot` and paste it to `frontend/languages/locale` directory.
2. Rename it with your language locale code e.g: `en_AU.po`, `de_DE.po`, etc.
3. Fill `msgstr` key with your translation. You can edit the `.po` file by using text editor or PO Edit software.
4. Pull request your translation to the `master` branch.

### Build Tania localisation by yourself

**THIS DOCUMENTATION WILL BE UPDATED LATER**

We are in the progress of building the new frontend application.

Then follow the instruction to [build Tania](#building-instructions).

## Support Us

We will move from OpenCollective to GitHub sponsorship. Thank you for all your donation in OpenCollective.

### Contributors

This project exists thanks to all the people who contribute.
<a href="https://github.com/tanibox/tania-core/graphs/contributors"><img src="https://opencollective.com/tania/contributors.svg?width=890&button=false" /></a>

### Backers

<a href="https://opencollective.com/tania"><img src="https://opencollective.com/tania/backers.svg?width=890&button=false" alt="backers"><img src="https://opencollective.com/tania/tiers/backer.svg?avatarHeight=36&width=600" alt="backers"></a>

## Copyright and License

Copyright to Tania and other contributors under [Apache 2.0](https://github.com/usetania/tania-core/blob/master/LICENSE) open source license.

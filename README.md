# Blog Go API 🚀

This repository contains the codes of the [Blog Go API] development by [phumiphatauk].

## Setup local development 🛠️

### Install tools 🧰

- [Docker desktop](https://www.docker.com/products/docker-desktop) 🐳
- [TablePlus](https://tableplus.com/) 🗄️
- [Golang](https://golang.org/) 🐹
- [Homebrew](https://brew.sh/) 🍺
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) 🔄

    ```bash
    brew install golang-migrate
    ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation) 🛠️

    ```bash
    brew install sqlc
    ```

### Setup infrastructure 🏗️

- Start postgres container:

    ```bash
    make postgres
    ```

- Create blog database:

    ```bash
    make createdb
    ```

- Run db migration up all versions:

    ```bash
    make migrateup
    ```

- Run db migration up 1 version:

    ```bash
    make migrateup no=1
    ```

- Run db migration down all versions:

    ```bash
    make migratedown
    ```

- Run db migration down 1 version:

    ```bash
    make migratedown_no no=1
    ```

### How to generate code 🧑‍💻

- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

- Create a new db migration:

    ```bash
    make migrate name=<migration_name>
    ```

### How to run 🚀

- Run server:

    ```bash
    make run
    ```


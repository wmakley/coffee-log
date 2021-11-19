# Coffee Log Go

Work in progress web app for logging your daily cup of coffee.

## Setup

### Dependencies

The following external tools are required, visit their websites for installation instructions:

* sqlc - https://sqlc.dev/
* dbmate - https://github.com/amacneil/dbmate

### Database

Create a ".env" file, with the following values:

```text
DATABASE_URL=postgres://[URL here]
TEST_DATABASE_URL=postgres://[URL here]
PORT=8080
```

Create those databases if coffee-log-go will not have permissions to do so, then run:

```shell
make migrate
```

### Compile and run unit tests

```shell
make
make test
```

# Coffee Log Go

Work in progress web app for logging your daily cup of coffee.

## Setup

Create a ".env" file, with the following values:

```text
DATABASE_URL=postgres://[URL here]
TEST_DATABASE_URL=postgres://[URL here]
PORT=8080
```

Create those databases if coffee-log-go will not have permissions to do so.

Then run:

```shell
make
make test
```

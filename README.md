# CQRS monitored App

Management books application (with volatile storage implementation) just to show how to apply best practices with CQRS, DDD, and EDA.

## Requirements

- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/) or [Orbstack](https://orbstack.dev/download)
- [Just](https://github.com/casey/just#installation)
- [Air](https://github.com/air-verse/air)
- [fzf](https://github.com/junegunn/fzf)
- [jq](https://jqlang.org/download/)
- [golang-ci lint](https://github.com/golangci/golangci-lint)

## API

```bash
curl -X POST \
                http://localhost:8080/api/books \
                -H 'content-type: application/json' \
                -d '{"id":"01JNVT8EENYPFTHJ5F4ZR1SN4E","title":"El Quijote","authorID":"01JNXWG511FY7YAC1QXM44BA9J"}'

curl -X PUT \
                http://localhost:8080/api/books/01JNVT8EENYPFTHJ5F4ZR1SN4E \
                -H 'content-type: application/json' \
                -d '{"id":"01JNVT8EENYPFTHJ5F4ZR1SN4E","title":"El ingenioso hidalgo Don Quijote de la Mancha","authorID":"01JNXWG511FY7YAC1QXM44BA9J"}'

curl -X GET http://localhost:8080/api/books | jq .

curl -X GET http://localhost:8080/api/books/01JNVT8EENYPFTHJ5F4ZR1SN4E | jq .

curl -X DELETE http://localhost:8080/api/books/01JNVT8EENYPFTHJ5F4ZR1SN4E
```

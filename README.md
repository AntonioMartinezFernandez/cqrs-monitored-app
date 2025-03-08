## API

```bash
curl -X POST \
                http://localhost:8000/api/books \
                -H 'content-type: application/json' \
                -d '{"id":"01JNVT8EENYPFTHJ5F4ZR1SN4E","title":"El ingenioso hidalgo Don Quijote de la Mancha","author":"Miguel de Cervantes"}'

curl -X GET http://localhost:8000/api/books | jq .

curl -X GET http://localhost:8000/api/books/01JNVT8EENYPFTHJ5F4ZR1SN4E | jq .
```

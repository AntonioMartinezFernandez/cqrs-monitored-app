## API

```bash
curl -X POST \
                http://localhost:8000/api/books \
                -H 'content-type: application/json' \
                -d '{"id":"01JNVT8EENYPFTHJ5F4ZR1SN4E","title":"el quijote","author":"miguel de cervantes"}'

curl -X GET http://localhost:8000/api/books
```

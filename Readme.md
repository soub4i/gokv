## Yay another concurrent kv store with golang

### How to run

```bash
# Run server
PORT=8080 go run main.go
```

### API

```bash
# Set key
curl GET http://localhost:8080/set/ke/value

# Get key

curl GET http://localhost:8080/get/key

# Delete key

curl GET http://localhost:8080/delete/key
```

Cheers!

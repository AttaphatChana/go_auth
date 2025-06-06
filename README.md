To test the service
```azure
docker compose build --no-cache                                             
docker compose up -d
```
test signup and login
```azure
curl -s -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"Secret123!"}' \
  | jq .
  
 curl -s -X POST http://localhost:8080/login \                               
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"Secret123!"}' \
  | jq .

```


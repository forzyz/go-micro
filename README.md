To Generate Protobuff command
```
protoc --go_out=./pb --go-grpc_out=./pb account.proto 
```

To docker build
```
docker compose up --build (or docker-compose up --build if you using legacy version)
```
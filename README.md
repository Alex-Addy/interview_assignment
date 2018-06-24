In order to run and test the server the following commands will be useful:

 - Run server: `go run main.go`
 - Send hash request: `curl --data "password=angrymonkey" http://localhost:8080/hash`
 - Send shutdown request: `curl http://localhost:8080/stop`

To successfully run automated tests make sure nothing else is currently binding `localhost:8080` and then type `go test`.


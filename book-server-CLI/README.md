* The following command will run the book server at default port 8081
```
    go run main.go
```
* we can set the port number using these flag
```
    go run main.go --port= 8081
    or,
    go run main.go -p= 8081
```
* These don't need authentication
```
    go run main.go -port=8081 --logIn=false
```

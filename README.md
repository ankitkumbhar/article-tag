# article-tag


- **article-tag** is a demo repo representing the use of dynamodb with golang. 
- With this example you will get the deep understanding of dynamodb, factors to keep while designing dynamodb.

## Setup process

### Running Application
To build and run application locally:

```shell
make run
```

### Testing
Used `testing` package that is built-in in Golang. To run unit tests run following command

```shell
go test -v ./... -cover -coverprofile=coverage.txt
```

To check the coverage run
```shell
go tool cover -func coverage.txt
```
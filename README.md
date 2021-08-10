# access-logs-parser

This is a simple golang-app, that takes [nginx-access-logs](https://docs.nginx.com/nginx/admin-guide/monitoring/logging/) 
as stdin and transforms it to csv-format into stdout.

## Running

```shell
cat access.log | go run . > result.csv
```


## Building

Following command builds this project into executable binary file: 
```shell
go build .
```

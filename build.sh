export CGO_ENABLED=0
export GOOS=linux
export GOHOSTOS=linux
export GOARCH=amd64
go build -o blogger ./main
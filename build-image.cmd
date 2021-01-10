set GOOS=linux
call go mod download
call go build -o bin/test-avito-mx cmd/test/main.go
call go build -o bin/simulator cmd/simulator/main.go
call docker build -t docker.io/elkozlova/test-avito-mx:latest .
call docker push docker.io/elkozlova/test-avito-mx:latest

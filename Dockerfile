FROM golang:latest 
 
WORKDIR /app 

COPY ./ /app

RUN go mod download && go build -o main .

EXPOSE 8080 50051

ENTRYPOINT go run main.go
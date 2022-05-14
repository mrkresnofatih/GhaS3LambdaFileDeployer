FROM golang:1.17

WORKDIR /app

COPY . .

RUN go build -o GhaS3KambdaFileDeployer

CMD ./GhaS3KambdaFileDeployer
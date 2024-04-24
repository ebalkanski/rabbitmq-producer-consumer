FROM golang:1.22.1

RUN go install github.com/ysmood/kit/cmd/guard@v0.25.11

WORKDIR /consumer

COPY . .

RUN go mod download

ENV TZ=Europe/Sofia

EXPOSE 81

ENTRYPOINT ["sh", "-c", "/go/bin/guard -w '**/*.go' -- go run ./cmd/..."]

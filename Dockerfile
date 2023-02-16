FROM golang:1.19-alpine as builder

WORKDIR /app

COPY ./main.go ./
COPY ./app ./app
RUN apk update && apk add --no-cache git
RUN GOARCH=amd64 GOOS=linux go mod init memePage && go mod tidy -v && go build -o main .


FROM alpine:latest

WORKDIR /usr/src/app

COPY --from=builder /app/main .

EXPOSE 6655

CMD ["./main"]

# docker build --platform linux/amd64 -t nesqui/meme-go .

#docker run -d --name meme-site --restart always \
#-v /opt/content:/usr/src/app/content --network host -p 6655:6655 --env-file dot.env nesqui/meme-go:latest

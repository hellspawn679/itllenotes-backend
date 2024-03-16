
FROM golang:1.21
WORKDIR /home/nekonotes
COPY . .
EXPOSE 7000
RUN go mod tidy 
RUN go build main.go
CMD ./main
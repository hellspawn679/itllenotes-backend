
FROM golang:1.21
WORKDIR /home/nekonotes
COPY . .
RUN go mod tidy 
RUN go build main.go
EXPOSE 7000
CMD ./main
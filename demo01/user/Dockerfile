FROM golang:1.15
WORKDIR /home/rookie/code/DStorage/user
COPY / /home/rookie/code/DStorage/user
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o user
ENTRYPOINT ["./user"]

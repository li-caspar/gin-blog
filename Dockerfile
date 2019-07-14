FROM alpine

WORKDIR $GOPATH/src/caspar/gin-blog
COPY . $GOPATH/src/caspar/gin-blog

EXPOSE 8080
ENTRYPOINT ["./gin-blog"]
FROM alpine:3.8

# HTTPSでのAPIリクエストに必要な証明書類
RUN apk add --no-cache ca-certificates

COPY bgslogviewer app/
COPY template app/template/
COPY main.html app/
COPY static app/static/

# Windowsでビルドした時用に、パーミッションを設定
RUN chmod -R 644 app
RUN chmod 744 app/bgslogviewer

WORKDIR app

CMD ["./bgslogviewer"]

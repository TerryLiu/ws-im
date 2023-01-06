FROM alpine:3.16

RUN apk add --update ca-certificates
RUN update-ca-certificates
RUN apk add --update tzdata
ENV TZ=Asia/Shanghai
COPY wsim /
RUN mkdir -p /usr/local/go/lib/time/
COPY zoneinfo.zip /usr/local/go/lib/time/
RUN chmod +x /wsim

CMD ["/wsim"]

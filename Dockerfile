FROM alpine:latest
COPY ./bin/report-generator .
EXPOSE 80
ENTRYPOINT ["./report-generator"]
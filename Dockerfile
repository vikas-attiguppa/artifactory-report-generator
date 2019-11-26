FROM alpine:latest
COPY ./bin/report-generator .
EXPOSE 8080
ENTRYPOINT ["./report-generator"]
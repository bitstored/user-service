FROM golang:alpine as source
WORKDIR /home/server
COPY . .
WORKDIR cmd/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -mod vendor -o user-service

FROM alpine as runner
LABEL name="bitstored/user-service"
RUN apk --update add ca-certificates
COPY --from=source /home/server/cmd/user-service/user-service /home/user-service
COPY --from=source /home/server/scripts/server.* /home/scripts/

WORKDIR /home
EXPOSE 4008
EXPOSE 5008

ENTRYPOINT [ "./user-service" ]

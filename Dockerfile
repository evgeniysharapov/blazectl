# Building 

FROM golang:1.21-alpine3.19 as builder
WORKDIR /blazectl
COPY . .
# certs business
RUN apk --no-cache add ca-certificates
# this is only needed for CISCO/Zscaler Umbrella Root CA
# building image from behind CISCO VPN
COPY cisco_umbrella_root.crt /etc/ssl/certs/

RUN go build -o /go/bin/blazectl

# Deployment 
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /go/bin/blazectl /app/blazectl
ENTRYPOINT [ "/app/blazectl" ]
CMD [ "help" ]
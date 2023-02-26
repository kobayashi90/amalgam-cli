FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN make _build-alpine


FROM alpine:latest AS runner
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/adcl .
RUN chmod +x /root/adcl
ENTRYPOINT ["/root/adcl"]
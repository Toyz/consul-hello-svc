FROM golang:latest as builder
WORKDIR /app
ARG project
COPY . .
RUN CGO_ENABLED=0 go build -o out.bin hello-world-svc

FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/out.bin /out
CMD ["/out"]
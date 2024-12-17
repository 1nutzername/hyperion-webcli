# Build stage - Build binary file
FROM golang:1.23.4-bookworm as builder

WORKDIR /build

COPY . .

RUN go get

# CGO_ENABLED=0 because I don't want to mess with C libraries. 
# -tags timetzdata to embed timezone database in case base image does not have.
# -trimpath to support reproduce build.
# -ldflags="-s -w" to strip debugging information.
RUN CGO_ENABLED=0 go build -o ./app -tags timetzdata -trimpath -ldflags="-s -w" .

# Run stage
FROM gcr.io/distroless/static-debian12

# Copy excecutable to distroless image
COPY --from=builder /build/app /app
COPY --from=builder /build/templates /templates

EXPOSE 8040

ENTRYPOINT [ "/app" ]

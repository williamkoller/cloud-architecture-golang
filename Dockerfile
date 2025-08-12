FROM golang:1.24.5-alpine AS deps
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

FROM deps AS builder

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG CGO_ENABLED=0

RUN CGO_ENABLED=$CGO_ENABLED GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o /out/bootstrap ./app/main.go


FROM public.ecr.aws/lambda/provided:al2

COPY --from=builder /out/bootstrap ${LAMBDA_RUNTIME_DIR}/bootstrap

RUN chmod +x ${LAMBDA_RUNTIME_DIR}/bootstrap

CMD ["bootstrap"]

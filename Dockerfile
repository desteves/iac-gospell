# syntax=docker/dockerfile:1

# build
FROM golang AS build-test-stage
WORKDIR /app
COPY . ./

# run integration tests
RUN --mount=type=secret,id=PULUMI_ACCESS_TOKEN,required=true \
    --mount=type=secret,id=GOOGLE_CREDENTIALS,required=true \
    export PULUMI_ACCESS_TOKEN=$(cat /run/secrets/PULUMI_ACCESS_TOKEN) && \
    export GOOGLE_CREDENTIALS=$(cat /run/secrets/GOOGLE_CREDENTIALS) && \
    curl -fsSL https://get.pulumi.com | sh && \
    export PATH=$PATH:/root/.pulumi/bin:/go/bin && \
    go mod download -x && \
    go install github.com/onsi/ginkgo/v2/ginkgo && \
    ginkgo run dev

# create binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /gospell-deploy

# Create Stack 
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-test-stage /gospell-deploy /gospell-deploy
ENTRYPOINT ["/gospell-deploy"]

# Build astro
FROM node:24.21.0-alpine AS web
RUN npm install --global pnpm@12.5.1
WORKDIR /src/web
COPY web/package.json web/pnpm-lock.yaml web/pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm build

# Build go binary
FROM golang:1.27.1-alpine AS server
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
# CGO_ENABLED=0 = static binary so it runs on distroless
# -ldflags="-s -w": discard symbol table and debug info for a smaller binary
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/gordle ./cmd/gordle

# assemble into distroless image with unpriviledged user
FROM gcr.io/distroless/static-debian13:nonroot
WORKDIR /app
COPY --from=server /out/gordle ./gordle
COPY --from=web /src/web/dist ./web/dist
ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT ["/app/gordle"]

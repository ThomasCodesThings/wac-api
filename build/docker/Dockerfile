# use specific versions of images
FROM openapitools/openapi-generator-cli:v7.0.1 as api

WORKDIR /local

COPY . .

# not used normally but redefine entrypoint for the case of checking this stage results
ENTRYPOINT ["bash"]

############################################

FROM golang:latest AS build

WORKDIR /app

# download dependencies - low frequency of changes
COPY go.mod go.sum ./
RUN go mod download

# copy sources - higher frequency of changes
COPY . .

RUN ls -la /app/internal

# ensure tests are passing
# RUN go test ./...

# create executable - department-webapi-srv
# we want to use scratch image so setting
# the build options in the way that will link all dependencies statically
RUN CGO_ENABLED=0 GOOS=linux \
    go build \
      -ldflags="-w -s" \
      -installsuffix 'static' \
      -o ./department-webapi-srv ./cmd/department-api-service

############################################

FROM scratch

# see https://github.com/opencontainers/image-spec/blob/main/annotations.md for details
LABEL org.opencontainers.image.authors="Tomas Cernacek"
LABEL org.opencontainers.image.title="Department WebAPI Service"
LABEL org.opencontainers.image.description="WEBAPI for Department"

# list all variables and their default values for clarity
ENV DEPARTMENT_API_ENVIRONMENT=production
ENV DEPARTMENT_API_PORT=8080
ENV DEPARTMENT_API_MONGODB_HOST=mongo
ENV DEPARTMENT_API_MONGODB_PORT=27017
ENV DEPARTMENT_API_MONGODB_DATABASE=cernica-department
ENV DEPARTMENT_API_MONGODB_COLLECTION=operation
ENV DEPARTMENT_API_MONGODB_USERNAME=rootf
ENV DEPARTMENT_API_MONGODB_PASSWORD=neUhaDnes
ENV DEPARTMENT_API_MONGODB_TIMEOUT_SECONDS=5

COPY --from=build /app/department-webapi-srv ./

# Actual port may be changed during runtime
# Default using for the simple case scenario
EXPOSE 8080
ENTRYPOINT ["./department-webapi-srv"]

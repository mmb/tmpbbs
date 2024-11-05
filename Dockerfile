ARG UID="1000"

FROM golang:1.23 AS build
ARG UID
RUN useradd --uid ${UID} tmpbbs
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags "-s -w" -v

FROM scratch
ARG UID
COPY --from=build /app/tmpbbs /tmpbbs
EXPOSE 8080/tcp
USER ${UID}
ENTRYPOINT ["/tmpbbs"]

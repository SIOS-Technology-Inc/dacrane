FROM docker:24.0.7-cli-alpine3.18

WORKDIR /build

# Install PDK
COPY ./src/pdk /usr/lib/go/src/dacrane/pdk

# Copy Source Code
COPY ./src/plugin/debug/ ./

# Install Go
RUN apk add --update --no-cache vim git make musl-dev go curl
ENV GOPATH=/root/go
ENV PATH=${GOPATH}/bin:/usr/local/go/bin:$PATH

# Install Plugin
RUN go build -o /usr/local/bin/plugin main.go
RUN chmod +x /usr/local/bin/plugin

WORKDIR /work

ENTRYPOINT [ "/usr/local/bin/plugin" ]

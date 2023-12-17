######################
#    Chrome Driver   #
######################

# Github: https://github.com/GoogleChromeLabs/chrome-for-testing
# Releases: https://googlechromelabs.github.io/chrome-for-testing/

FROM debian:stable-slim as chrome

# Never ask for user input
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get upgrade --yes && \
    apt-get install --yes \
        ca-certificates \
        openssl \
        curl \
        unzip \
        jq && \
    export CHROME_VERSION=$(curl -s https://googlechromelabs.github.io/chrome-for-testing/last-known-good-versions.json | jq ."channels"."Stable"."version" | sed 's/"//g') && \
    curl -o /opt/chromedriver.zip https://edgedl.me.gvt1.com/edgedl/chrome/chrome-for-testing/${CHROME_VERSION}/linux64/chromedriver-linux64.zip && \
    curl -o /opt/chrome.zip https://edgedl.me.gvt1.com/edgedl/chrome/chrome-for-testing/${CHROME_VERSION}/linux64/chrome-linux64.zip && \
    cd /opt && \
    unzip /opt/chromedriver.zip && \
    unzip /opt/chrome.zip && \
    mv /opt/chromedriver-linux64/chromedriver /chromedriver
    # TBD move chrome to root

######################
#      Builder       #
######################

FROM golang:1.20 as builder

# Never ask for user input
ARG DEBIAN_FRONTEND=noninteractive

# ACCEPT_EULA=Y is required to install Microsoft ODBC Driver
ARG ACCEPT_EULA=Y

RUN apt-get update && \
    apt-get upgrade --yes && \
    apt-get install --yes \
        ca-certificates \
        openssl \
        curl \
        gnupg2 && \
    curl https://packages.microsoft.com/keys/microsoft.asc | tee /etc/apt/trusted.gpg.d/microsoft.asc && \
    gpg --dearmor < /etc/apt/trusted.gpg.d/microsoft.asc > /usr/share/keyrings/microsoft-prod.gpg && \
    curl https://packages.microsoft.com/config/debian/12/prod.list | tee /etc/apt/sources.list.d/mssql-release.list && \
    apt-get update && \
    apt-get install --yes \
        unixodbc \
        unixodbc-dev \
        # install Microsoft ODBC Driver for SQL Server
        msodbcsql18 \
        # optional: for bcp and sqlcmd
        mssql-tools18

# Set destination for COPY
WORKDIR /build

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY cmd ./cmd/
COPY internal ./internal/

# RUN ls -alh /build/

# Build
RUN GOARCH=amd64 GOOS=linux go build -o /wb-parsing-crawler cmd/main.go

######################
#       Runner       #
######################

FROM debian:stable-slim as runner

ARG USERNAME=crawler
ARG WORK_DIR=/opt/crawler

# Never ask for user input
ARG DEBIAN_FRONTEND=noninteractive

# ACCEPT_EULA=Y is required to install Microsoft ODBC Driver
ARG ACCEPT_EULA=Y

ENV CONFIG_TYPE yaml
ENV CONFIG_PATH /usr/local/etc/crawler.yaml

COPY --from=chrome /chromedriver /usr/local/bin/chromedriver
COPY --from=builder /wb-parsing-crawler /usr/local/bin/wb-parsing-crawler

RUN apt-get update && \
    apt-get upgrade --yes && \
    apt-get install --yes \
        ca-certificates \
        openssl \
        curl \
        gnupg2 \
        libglib2.0-0 \
        libnss3 \
        libxcb1 && \
    curl https://packages.microsoft.com/keys/microsoft.asc | tee /etc/apt/trusted.gpg.d/microsoft.asc && \
    gpg --dearmor < /etc/apt/trusted.gpg.d/microsoft.asc > /usr/share/keyrings/microsoft-prod.gpg && \
    curl https://packages.microsoft.com/config/debian/12/prod.list | tee /etc/apt/sources.list.d/mssql-release.list && \
    apt-get update && \
    apt-get install --yes \
        unixodbc \
        # install Microsoft ODBC Driver for SQL Server
        msodbcsql18 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    adduser --no-create-home --disabled-login --shell /bin/nologin ${USERNAME} && \
    mkdir -p ${WORK_DIR} && \
    chown ${USERNAME}:${USERNAME} ${WORK_DIR} && \
    ln -s /usr/local/bin/chromedriver ${WORK_DIR}/chromedriver

WORKDIR ${WORK_DIR}/

USER ${USERNAME}
CMD ["/usr/local/bin/wb-parsing-crawler"]
# CMD ["-config", "$CONFIG_PATH", "-config-type", "$CONFIG_TYPE"]

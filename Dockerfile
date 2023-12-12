FROM debian:stable-slim

ARG USERNAME=crawler
ARG WORK_DIR=/opt/crawler

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
        mssql-tools18 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    adduser --no-create-home --disabled-login --shell /bin/nologin ${USERNAME} && \
    mkdir -p ${WORK_DIR}

WORKDIR ${WORK_DIR}/

# USER ${USERNAME}
# CMD ["main.py"]

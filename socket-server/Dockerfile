FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    curl \
    wget \
    git \
    vim \
    build-essential \
    software-properties-common \
    ca-certificates \
    apt-transport-https \
    lsb-release \
    gnupg \
    && rm -rf /var/lib/apt/lists/*

# Non root user.
RUN useradd -ms /bin/bash appuser

#Go setup
RUN wget https://go.dev/dl/go1.22.1.linux-amd64.tar.gz -O /tmp/go.tar.gz && \
    tar -C /usr/local -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz

ENV PATH="/usr/local/go/bin:${PATH}"

# Node and nvm setup.
ENV NVM_DIR="/usr/local/nvm"
ENV NODE_VERSION="18.12.0"
RUN mkdir $NVM_DIR

RUN curl https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

ENV NODE_PATH=$NVM_DIR/v$NODE_VERSION/lib/node_modules
ENV PATH=$NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

RUN echo "source $NVM_DIR/nvm.sh && \
    nvm install $NODE_VERSION && \
    nvm alias default $NODE_VERSION && \
    nvm use default" | bash

RUN curl -qL https://www.npmjs.com/install.sh | sh

RUN apt-get update && apt-get install -y python3 python3-pip && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o socket-server

#Expose 2 ports for user to use.
EXPOSE 5000

EXPOSE 3000

CMD ["./socket-server"]

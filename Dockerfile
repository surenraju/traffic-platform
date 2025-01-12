# Use devopsinfra/docker-terragrunt as the base image
FROM devopsinfra/docker-terragrunt:aws-latest

# Install required packages and update Golang
RUN apt-get update && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:longsleep/golang-backports -y && \
    apt-get update && \
    apt-get install -y golang-go && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set the updated Go binary path
ENV PATH="/usr/lib/go/bin:${PATH}"

# Verify installation
RUN go version
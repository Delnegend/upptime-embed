FROM debian:12-slim

RUN apt-get update && apt-get install -y \
    curl git unzip sudo && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

RUN groupadd -g 1000 upptime && \
    useradd -m -u 1000 -g upptime upptime && \
    echo "upptime ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/upptime

USER upptime
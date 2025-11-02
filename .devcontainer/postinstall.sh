#!/usr/bin/env bash

# go
sudo rm -rf /usr/local/go && curl -fsSL https://go.dev/dl/go1.25.3.linux-amd64.tar.gz | sudo tar -C /usr/local -xzf -
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
export PATH=$PATH:/usr/local/go/bin

# js
curl -fsSL https://bun.sh/install | bash
echo 'export BUN_INSTALL="$HOME/.bun"' >> ~/.bashrc
echo 'export PATH="$BUN_INSTALL/bin:$PATH"' >> ~/.bashrc

# just
ver="1.43.0"
curl -L https://github.com/casey/just/releases/download/$ver/just-$ver-x86_64-unknown-linux-musl.tar.gz | sudo tar -xzf - -C /usr/local/bin just
echo 'alias j=just' >> ~/.bashrc
echo 'eval "$(just --completions bash)"' >> ~/.bashrc
echo 'complete -F _just j' >> ~/.bashrc

# fzf
ver="0.66.0"
curl -L https://github.com/junegunn/fzf/releases/download/v$ver/fzf-$ver-linux_amd64.tar.gz | sudo tar -xzf - -C /usr/local/bin fzf
FROM golang:1.24

RUN apt-get update && \
    apt-get install -y zsh tree && \
    rm -rf /var/lib/apt/lists/*

RUN useradd -m vscode
USER vscode
WORKDIR /workspace

ENV PATH=$PATH:/usr/local/go/bin:/go/bin

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.0.2

RUN sh -c "$(wget https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh -O -)" && \
    git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions && \
    git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting && \
    git clone https://github.com/MichaelAquilina/zsh-you-should-use.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/you-should-use

COPY --chown=vscode:vscode .zshrc /home/vscode/.zshrc

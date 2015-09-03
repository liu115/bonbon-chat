FROM base/archlinux:latest

COPY bonbon-*.tar.gz .
COPY mirrorlist /etc/pacman.d/mirrorlist

RUN pacman -Sy --noconfirm
RUN pacman -S --noconfirm go npm git gcc
RUN npm install -g jake
RUN tar xf bonbon-*.tar.gz; \
    rm bonbon-*.tar.gz; \
    cd bonbon-*; \
    mkdir .go; \
    export GOPATH="$PWD/.go"; \
    export PATH="$PWD/.go/bin:$PATH"; \
    go get github.com/mattn/gom; \
    jake clean; \
    jake

EXPOSE 8080
CMD ["/bin/bash", "-c", "cd /bonbon-*; jake test > /var/log/bonbon.log && exec ./bonbon-server >> /var/log/bonbon.log 2>&1"]

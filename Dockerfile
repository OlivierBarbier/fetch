FROM        golang
RUN         go get github.com/OlivierBarbier/fetch
ENTRYPOINT  ["fetch"]
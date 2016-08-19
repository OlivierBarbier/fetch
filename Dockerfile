FROM        golang
RUN         go get github.com/OlivierBarbier/fetch
EXPOSE      8080
ENTRYPOINT  ["fetch"]
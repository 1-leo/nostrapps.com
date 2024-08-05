dev:
    fd 'go|html|md|templ|base.css' | entr -r bash -c 'just build && godotenv ./nostrapps.com'

templ:
    templ generate

build:
    just templ
    CC=musl-gcc go build -ldflags="-s -w -linkmode external -extldflags '-static' -s -w" -o ./

deploy: build
    ssh root@erhard 'systemctl stop nostrapps.com'
    scp nostrapps.com erhard:nostrapps.com/nostrapps.com
    ssh root@erhard 'systemctl start nostrapps.com'

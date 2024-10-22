export PATH := "./node_modules/.bin:" + env_var('PATH')

dev:
    fd 'go|html|md|templ|base.css' | entr -r bash -c 'just build && godotenv ./nostrapps.com'

build: templ tailwind
    CC=musl-gcc go build -ldflags="-s -w -linkmode external -extldflags '-static' -s -w" -o ./

templ:
    templ generate

tailwind:
    tailwind -i base.css -o static/bundle.css

deploy target: build
    ssh root@{{target}} 'systemctl stop nostrapps.com'
    scp nostrapps.com {{target}}:nostrapps.com/nostrapps.com
    ssh root@{{target}} 'systemctl start nostrapps.com'

#!/usr/bin/env bash

lint() {
  golangci-lint run ./...
}

build() {
  go build -o bin/simplesrv .
}

test() {
  go test -test.count=1 -cover ./...
}

coverage() {
  local tmp=/tmp/simplesrv-profile.out
  go test -coverprofile="${tmp}" ./... >/dev/null
  go tool cover -func="${tmp}"
  rm -f "${tmp}"
}

run() {
  ./bin/simplesrv
}

help() {
  printf "%s <task> <args>\n" "$0"
  printf "Tasks:\n"
  compgen -A function | cat -n
}

die() {
  printf "%s\n" "$@"
  exit 1
}

action="$1"
case $action in
  lint | build | build-server | run | test | coverage | help)
    "$@"
    ;;
  *)
    die "invalid action '${action}'"
    ;;
esac

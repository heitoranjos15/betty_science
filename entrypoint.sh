#!/bin/sh

if [ "${DLV_DEBUG}" = "true" ]; then
    dlv debug --listen=:2345 --headless=true --api-version=2 --output=/.dlv/__debug_bin --log "$*" && exit 1
else
    reflex -s -r '\.(go|tpl|tmpl)$' -R '^vendor/' -- sh -c "go run -buildvcs=false $*"
fi


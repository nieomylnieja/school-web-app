#!/bin/sh

if [ -z "$BACKEND_PORT" ]; then
	echo "BACKEND_PORT not defined"
	exit 1
fi

sed -i "s~9000~${BACKEND_PORT}~g" ./static/js/http.js

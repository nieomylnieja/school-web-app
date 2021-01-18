#!/bin/sh

if [ -z "$BACKEND_URL" ]; then
	echo "BACKEND_URL not defined"
	exit 1
fi

sed -i "s~localhost:3000~${BACKEND_URL}~g" ./static/js/http.js

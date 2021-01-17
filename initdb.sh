#!/bin/bash
set -Eeuo pipefail

if [ "$MONGO_USERNAME" ] && [ "$MONGO_PASSWORD" ]; then
  "${mongo[@]}" -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" --authenticationDatabase "$rootAuthDatabase" "$MONGO_DATABASE" <<-EOJS
        db.createUser({
            user: $(_js_escape "$MONGO_USERNAME"),
            pwd: $(_js_escape "$MONGO_PASSWORD"),
            roles: [ { role: 'readWrite', db: $(_js_escape "$MONGO_DATABASE") } ]
        })
        db = new Mongo().getDB("$MONGO_DATABASE");
        db.createCollection('user', { capped: false });
        db.createCollection('password', { capped: false });
        db.user.createIndex({'email': 1},{unique: true, name: 'uniqueEmail'})
        db.password.createIndex({'userId': 1},{unique: true, name: 'uniqueUserId'})
EOJS
fi

#!/usr/bin/env bash

MONGO_HOST="localhost"
MONGO_PORT="27017"
MONGO_USER="root"
MONGO_PASS="userdb"
AUTH_DB="admin"
DB_NAME="userdb"

mongosh "mongodb://$MONGO_USER:$MONGO_PASS@$MONGO_HOST:$MONGO_PORT/$DB_NAME?authSource=$AUTH_DB"

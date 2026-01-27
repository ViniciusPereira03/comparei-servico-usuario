#!/usr/bin/env bash

MONGO_HOST="localhost"
MONGO_PORT="27017"
MONGO_USER="root"
MONGO_PASS="userdb"
AUTH_DB="admin"
DB_NAME="userdb"

docker exec -it mongo-user mongosh -u root -p userdb --authenticationDatabase admin

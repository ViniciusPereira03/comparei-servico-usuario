db = db.getSiblingDB('userdb');

db.users.createIndex({ "username": 1, "email": 1 }, { unique: true });

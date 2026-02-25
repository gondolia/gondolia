#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE catalog;
    GRANT ALL PRIVILEGES ON DATABASE catalog TO postgres;

    CREATE DATABASE cart;
    GRANT ALL PRIVILEGES ON DATABASE cart TO postgres;

    CREATE DATABASE "order";
    GRANT ALL PRIVILEGES ON DATABASE "order" TO postgres;
EOSQL

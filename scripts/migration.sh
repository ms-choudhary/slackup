#! /bin/bash

rm -f "$PWD/slackup.db"
sqlite3 slackup.db < "$PWD/scripts/migration.sql"

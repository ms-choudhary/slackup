#! /bin/bash

rm -f ./slackup.db
sqlite3 slackup.db < migration.sql

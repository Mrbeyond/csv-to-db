#!/bin/sh

#  Delete production folder is exist for new bundle
rm -rf production
# Make new dir for production fils
mkdir -p production

# Copy project folder exclude redundant local files

rsync -avP --delete --exclude .git --exclude Dockerfile.local --exclude .env \
--exclude main ./project ./production/

# copy docker compose and run bash
cp docker-compose.yml ./production/
cp run.sh ./production/
cp readme.md ./production/
cp  prod.env ./production/
mv ./production/prod.env ./production/.env

mkdir -p production/psql_data

# Rename project/ prod.env to .env
ls -al  ./production/project
mv ./production/project/prod.env ./production/project/.env
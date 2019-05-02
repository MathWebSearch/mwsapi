#!/bin/bash
set -e

docker-compose up --force-recreate -d mwsquery_mws

printf "Waiting for server to be up ..."
while [[ "$(curl -s -o /dev/null -w ''%{http_code}'' localhost:1000)" != "200" ]]; do 
  sleep 5
done

printf " done\e[0m\n"
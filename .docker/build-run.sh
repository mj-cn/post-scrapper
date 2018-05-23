#!/usr/bin/env bash

docker build -t post-scrapper ../
docker run --link mongo:mongo post-scrapper http://www.venusgo.com/search/node/clasex?page=1 mongo:27017
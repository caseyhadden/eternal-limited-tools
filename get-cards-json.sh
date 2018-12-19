#!/bin/sh

curl https://eternalwarcry.com/content/cards/eternal-cards.json | jq . > data/eternal-cards.json

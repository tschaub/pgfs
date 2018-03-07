# Postgres backed WFS 3

## Build it

    go build -o pgfs main.go

## Run it

    pgfs serve "dbname=pgfs sslmode=disable"

## Sample requests

### list all collections
    curl -s http://localhost:5000/collections --header "Accept: application/json" | jj -p

### create a new collection
    curl -s http://localhost:5000/collections \
      --header "Content-Type: application/json" \
      --data '{"name": "countries", "title": "Countries", "description": "Countries of the world"}' | jj -p

### post features to a collection
    curl -s http://localhost:5000/collections/countries/items \
      --request POST \
      --header "Content-Type: application/json" \
      --data @testdata/countries.json

### get features in a collection
    curl -s http://localhost:5000/collections/countries/items | jj -p

# Postgres Backed WFS 3

    pgfs serve "dbname=pgfs sslmode=disable"

## Sampe requests

# list all collections
curl -s http://localhost:5000/collections --header "Accept: application/json" | jj -p

# create a new collection
curl -s http://localhost:5000/collections \
  --header "Content-Type: application/json" \
  --data '{"name": "countries", "title": "Countries", "description": "Countries of the world"}' | jj -p

# post features to a collection
curl -s http://localhost:5000/collections/countries/items \
  --request POST \
  --header "Content-Type: application/json" \
  --data @testdata/countries.json

# get features in a collection
curl -s http://localhost:5000/collections/countries/items | jj -p

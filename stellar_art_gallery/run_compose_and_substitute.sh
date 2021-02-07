#!/bin/bash

docker-compose up -d
sleep 10
docker-compose exec -w /usr/share/nginx/html/static/js client_stellar_art_gallery bash -c "sed -i 's,__API_URL__,https://galleryapi.vatnar.no,g' *"


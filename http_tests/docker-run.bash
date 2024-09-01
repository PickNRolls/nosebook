cd ..
docker container stop db
docker container rm db
docker compose up --build --wait
cd -

docker build -t backend-http-tests .
docker run --rm \
  --network back_default \
  -it -e \"TERM=xterm-256color\" \
  -v ./src:/usr/src/app/src \
  backend-http-tests \
  npm run docker-start -- $@


docker run \
  -d \
  --rm \
  --name mongo4.4 \
  -p 27017:27017 \
  -v "$(pwd)/.dev/data/mongo":/data/db \
  mongo:4.4

echo "Stopping docker container..."
docker-compose down

echo "Building..."
mkdir build
env GOOS=linux GOARCH=amd64 go build -o build/api . 

echo "Starting docker container..."
docker-compose up
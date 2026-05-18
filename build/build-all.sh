version=2

docker build --platform linux/amd64 -t williamqm/mmobackend:latest .
docker push williamqm/mmobackend:v2

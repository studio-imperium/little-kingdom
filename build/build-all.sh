version=2

// build go binary
GOOS=linux GOARCH=amd64 go build -o backend

// add to dockrer
docker build --platform linux/amd64 -t williamqm/mmobackend:latest .
docker push williamqm/mmobackend:latest


docker compose down
rm -rf /root/Caddyfile
cat > /root/Caddyfile <<'EOF'
server.kingdomcrushers.com {
reverse_proxy backend:8082
}
EOF
docker compose up -d

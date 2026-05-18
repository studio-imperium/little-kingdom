FROM --platform=linux/amd64 ubuntu:latest

WORKDIR /workspace

COPY backend /workspace/backend
COPY restart_server.sh /workspace/restart_server.sh

EXPOSE 20200
RUN chmod +x /workspace/backend
RUN chmod +x /workspace/restart_server.sh
CMD ["/workspace/restart_server.sh"]

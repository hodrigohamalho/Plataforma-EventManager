FROM pmoneda/go_alpine:latest
ADD event-manager /
ENV PORT 8081
RUN chmod +x event-manager
ENTRYPOINT ["/event-manager"]
EXPOSE 8081
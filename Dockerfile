FROM cgr.dev/chainguard/static:latest

ARG UID=65532
ARG GID=65532

WORKDIR /app

COPY --chown=${UID}:${GID} hackerone-exporter /app/hackerone-exporter

EXPOSE 8080

USER nonroot

ENTRYPOINT [ "/app/hackerone-exporter" ]

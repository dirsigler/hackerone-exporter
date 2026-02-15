FROM cgr.dev/chainguard/busybox:latest

COPY hackerone-exporter /usr/bin/hackerone-exporter

EXPOSE 9194

USER nonroot:nonroot

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:9194/metrics || exit 1

ENTRYPOINT ["/usr/bin/hackerone-exporter"]

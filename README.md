# Prometheus Exporter for HackerOne

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/dirsigler/hackerone-exporter)](https://goreportcard.com/report/github.com/dirsigler/hackerone-exporter)

This is a Prometheus Exporter for the [HackerOne](https://hackerone.com) API. It allows you to monitor your HackerOne data in Prometheus and build dashboards in Grafana.

## ⚙️ Metrics

| Name                                | Labels                     | Description                                  |
| ----------------------------------- | -------------------------- | -------------------------------------------- |
| `hackerone_assets_total`            | `organization_id`          | Total number of HackerOne Assets             |
| `hackerone_reports_total`           | `organization_id`          | Total number of HackerOne Reports            |
| `hackerone_programs_total`          | `handle`, `state`          | Total number of HackerOne Programs           |
| `hackerone_invited_hackers_total`   | `organization_id`, `state` | Total number of HackerOne Invited Hackers    |
| `hackerone_weaknesses_total`        | `name`, `id`               | Total number of HackerOne Weaknesses         |
| `hackerone_scrape_errors_total`     |                            | Total number of HackerOne API scrape errors  |
| `hackerone_last_scrape_timestamp`   |                            | Unix timestamp of the last successful scrape |
| `hackerone_scrape_duration_seconds` |                            | Duration of HackerOne API scrapes in seconds |

## 🚀 Deployment

With each [release](https://github.com/dirsigler/hackerone-exporter/releases), a secure-by-default Docker image is available on [GitHub](https://github.com/dirsigler/hackerone-exporter/pkgs/container/hackerone-exporter) and [DockerHub](https://hub.docker.com/repository/docker/dirsigler/hackerone-exporter/general).

### Docker Compose

Here is a sample `docker-compose.yml`:

```yaml
version: "3.8"
services:
  hackerone-exporter:
    image: ghcr.io/dirsigler/hackerone-exporter:latest
    container_name: hackerone-exporter
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - HACKERONE_API_USER=<YOUR_API_USER>
      - HACKERONE_API_PASSWORD=<YOUR_API_PASSWORD>
      - HACKERONE_ORG_ID=<YOUR_ORG_ID>
```

### Docker

```sh
docker run --rm \
  --interactive --tty \
  --env HACKERONE_API_USER=<YOUR_API_USER> \
  --env HACKERONE_API_PASSWORD=<YOUR_API_PASSWORD> \
  --env HACKERONE_ORG_ID=<YOUR_ORG_ID> \
  ghcr.io/dirsigler/hackerone-exporter:latest
```

## 🚩 Configuration

`$ hackerone-exporter --help`

| Flag                | Environment Variable     | Description                          | Default                     |
| ------------------- | ------------------------ | ------------------------------------ | --------------------------- |
| `--api-user`        | `HACKERONE_API_USER`     | HackerOne API Username               | **required**                |
| `--api-password`    | `HACKERONE_API_PASSWORD` | HackerOne API Password               | **required**                |
| `--org-id`          | `HACKERONE_ORG_ID`       | HackerOne Organization ID            | **required**                |
| `--port`            | `PORT`                   | Port to listen on                    | `8080`                      |
| `--scrape-interval` | `SCRAPE_INTERVAL`        | Scrape interval in seconds           | `60`                        |
| `--log-level`       | `LOG_LEVEL`              | Log level (debug, info, warn, error) | `info`                      |
| `--api-url`         | `HACKERONE_API_URL`      | HackerOne API URL                    | `https://api.hackerone.com` |

## 📝 License

Built with ☕️ and licensed under the [Apache 2.0 License](./LICENSE).

# Prometheus Exporter for HackerOne

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/dirsigler/hackerone-exporter)](https://goreportcard.com/report/github.com/dirsigler/hackerone-exporter)

This is a Prometheus Exporter for the [HackerOne](https://hackerone.com) API. It allows you to monitor your HackerOne data in Prometheus and build dashboards in Grafana.

## ⚙️ Metrics

### Report Metrics

| Name                                           | Labels                             | Description                                  |
| ---------------------------------------------- | ---------------------------------- | -------------------------------------------- |
| `hackerone_reports_total_count`                | -                                  | Total number of reports across all programs  |
| `hackerone_reports_state_count`                | `state`                            | Number of reports by state                   |
| `hackerone_reports_severity_count`             | `severity`                         | Number of reports by severity                |
| `hackerone_reports_reporter_count`             | `reporter_username`, `reporter_id` | Number of reports submitted by each reporter |
| `hackerone_reports_time_to_triage_seconds`     | -                                  | Histogram of time from creation to triage    |
| `hackerone_reports_time_to_resolution_seconds` | -                                  | Histogram of time from creation to close     |

### Program Metrics

| Name                                                  | Labels       | Description                                    |
| ----------------------------------------------------- | ------------ | ---------------------------------------------- |
| `hackerone_programs_total_count`                      | -            | Total number of programs                       |
| `hackerone_programs_balance_usd`                      | `program_id` | Current bounty balance in USD                  |
| `hackerone_programs_transactions_total`               | `program_id` | Total number of payment transactions           |
| `hackerone_programs_transactions_bounty_usd_total`    | `program_id` | Total bounty amount paid in USD                |
| `hackerone_programs_allowed_reporters_total`          | `program_id` | Number of allowed reporters (private programs) |
| `hackerone_programs_allowed_reporters_verified_total` | `program_id` | Number of identity-verified reporters          |
| `hackerone_programs_swag_total`                       | `program_id` | Total number of swag items awarded             |
| `hackerone_programs_swag_sent_total`                  | `program_id` | Number of swag items sent                      |
| `hackerone_programs_swag_pending_total`               | `program_id` | Number of swag items pending shipment          |
| `hackerone_programs_triage_reviews_total`             | `program_id` | Total number of triage reviews received        |
| `hackerone_programs_triage_review_rating_average`     | `program_id` | Average triage review rating (1-5 scale)       |
| `hackerone_programs_common_responses_total`           | `program_id` | Number of common response templates            |
| `hackerone_programs_reporters_total`                  | `program_id` | Number of unique reporters who submitted       |
| `hackerone_programs_thanks_total`                     | `program_id` | Number of hackers in thanks list               |
| `hackerone_programs_weaknesses_total`                 | `program_id` | Number of weakness categories available        |
| `hackerone_programs_members_total`                    | `program_id` | Total number of program team members           |
| `hackerone_programs_admins_total`                     | `program_id` | Number of program administrators               |

### Exporter Metrics

| Name                                                           | Labels     | Description                                |
| -------------------------------------------------------------- | ---------- | ------------------------------------------ |
| `hackerone_exporter_up`                                        | -          | Whether the exporter is up (1) or down (0) |
| `hackerone_exporter_collection_duration_seconds`               | -          | Duration of last collection in seconds     |
| `hackerone_exporter_api_call_duration_seconds`                 | `endpoint` | Duration of API calls by endpoint          |
| `hackerone_exporter_api_call_errors_total`                     | `endpoint` | Total number of API call errors            |
| `hackerone_exporter_last_collection_success_timestamp_seconds` | -          | Timestamp of last successful collection    |
| `hackerone_exporter_last_collection_attempt_timestamp_seconds` | -          | Timestamp of last collection attempt       |
| `hackerone_exporter_collection_errors_total`                   | -          | Total number of collection errors          |

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
      - "9194:9194"
    environment:
      - HACKERONE_API_USER=<YOUR_API_USER>
      - HACKERONE_API_PASSWORD=<YOUR_API_PASSWORD>
      - HACKERONE_ORG_ID=<YOUR_ORG_ID>
```

### Docker

```sh
docker run --rm \
  --interactive --tty \
  --publish 9194:9194 \
  --env HACKERONE_API_USER=<YOUR_API_USER> \
  --env HACKERONE_API_PASSWORD=<YOUR_API_PASSWORD> \
  --env HACKERONE_ORG_ID=<YOUR_ORG_ID> \
  ghcr.io/dirsigler/hackerone-exporter:latest
```

## 🚩 Configuration

`$ hackerone-exporter --help`

| Flag             | Environment Variable        | Description                                                       | Default                     |
| ---------------- | --------------------------- | ----------------------------------------------------------------- | --------------------------- |
| `--api-username` | `HACKERONE_API_USER`        | HackerOne API Username                                            | **required**                |
| `--api-password` | `HACKERONE_API_PASSWORD`    | HackerOne API Password                                            | **required**                |
| `--org-id`       | `HACKERONE_ORG_ID`          | HackerOne Organization ID                                         | **required**                |
| `--port`         | `SERVER_PORT`               | Port to listen on                                                 | `9194`                      |
| `--log-level`    | `LOG_LEVEL`                 | Log level (debug, info, warn, error)                              | `info`                      |
| `--api-url`      | `HACKERONE_API_URL`         | HackerOne API URL                                                 | `https://api.hackerone.com` |
| `--programs`     | `HACKERONE_PROGRAM_HANDLES` | Specific program handles to collect metrics for (comma-separated) | all programs                |

### Examples

**Collect metrics for all programs (default):**

```sh
hackerone-exporter \
  --api-username=your-username \
  --api-password=your-password \
  --org-id=your-org-id
```

**Collect metrics for specific programs only:**

```sh
hackerone-exporter \
  --api-username=your-username \
  --api-password=your-password \
  --org-id=your-org-id \
  --programs=my-program-1 \
  --programs=my-program-2
```

**Or using environment variables:**

```sh
export HACKERONE_API_USER=your-username
export HACKERONE_API_PASSWORD=your-password
export HACKERONE_ORG_ID=your-org-id
export HACKERONE_PROGRAM_HANDLES=my-program-1,my-program-2

hackerone-exporter
```

## 📝 License

Built with ☕️ and licensed under the [Apache 2.0 License](./LICENSE).

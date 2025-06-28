# Prometheus Exporter for [HackerOne](https://HackerOne.com)

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/dirsigler/hackerone-exporter)](https://goreportcard.com/report/github.com/dirsigler/hackerone-exporter)

This is an custom Prometheus Exporter for the https://HackerOne.com featureflag solution.

Allows your business to monitor all feature-flags across your Organization, Product and Configs.
Also includes metrics for Zombie flags (stale feature-flags).
Improve your software application hygiene and better understand the usage of your [HackerOne.com](https://HackerOne.com) setup.

## ⚙️ Metrics

The HackerOne Prometheus Exporter supports all basic pre-configured types of incidents available in [HackerOne.com](https://HackerOne.com).

| Name                                | Label                                                    | Description                                        |
| ----------------------------------- | -------------------------------------------------------- | -------------------------------------------------- |
| `hackerone_products_total`          |                                                          | Total number of HackerOne products                 |
| `hackerone_configs_total`           | `product_name`, `product_id`                             | Total number of HackerOne configs per product      |
| `hackerone_environments_total_`     | `product_name`, `product_id`                             | Total number of HackerOne environments per product |
| `hackerone_feature_flags_total`     | `product_name`, `product_id`, `config_name`, `config_id` | Total number of feature flags per config           |
| `hackerone_scrape_errors_total`     |                                                          | Total number of HackerOne API scrape errors        |
| `hackerone_last_scrape_timestamp`   |                                                          | Unix timestamp of the last successful scrape       |
| `hackerone_scrape_duration_seconds` |                                                          | Duration of HackerOne API scrapes in seconds       |
| `hackerone_zombie_flags_total`      | `product_name`, `product_id`, `config_name`, `config_id` | Total number of zombie flags per product           |

## 🚀 Deployment

> IMPORTANT: You have to provide the "HackerOne_API_KEY="<MY_API_KEY>" environment variable to your deployment for the HackerOne Prometheus Exporter to work.

---

With each [release](https://github.com/dirsigler/hackerone-exporter/releases) I also provide a [secure by default](https://www.chainguard.dev/chainguard-images) Docker Image.

You can chose from:

- The Image on GitHub => [hackerone-exporter on GitHub](https://github.com/dirsigler/hackerone-exporter/pkgs/container/hackerone-exporter)
- The Image on DockerHub => [hackerone-exporter on DockerHub](https://hub.docker.com/repository/docker/dirsigler/hackerone-exporter/general)

### Docker

```sh
docker run --rm \
--interactive --tty \
--env HackerOne_API_KEY="<MY_API_KEY>" \
ghcr.io/dirsigler/hackerone-exporter:latest
```

You can also enable a logger with Debug mode via the `--log.level=debug` flag.
See the available [configuration](#🚩-configuration)

## 🚩 Configuration

```sh
$ hackerone-exporter --help

NAME:
   hackerone-exporter - Export HackerOne metrics to Prometheus

USAGE:
   hackerone-exporter [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --HackerOne-api-key value  HackerOne API key (can also be set via HackerOne_API_KEY env var) [$HACKERONE_API_KEY]
   --port value               Port to listen on (default: 8080) [$PORT]
   --scrape-interval value    Scrape interval in seconds (default: 60) [$SCRAPE_INTERVAL]
   --log-level value          Log level (debug, info, warn, error) (default: "info") [$LOG_LEVEL]
   --organization-id value    HackerOne Organization ID [$HACKERONE_ORGANIZATION_ID]
   --product-id value         HackerOne Product ID [$HACKERONE_PRODUCT_ID]
   --help, -h                 show help
```

## 📝 License

Built with ☕️ and licensed via [Apache 2.0](./LICENSE)

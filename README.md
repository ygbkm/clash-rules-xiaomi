# Clash rules for Xiaomi

[![Latest Release](https://img.shields.io/github/release-date/ygbkm/clash-rules-xiaomi?display_date=published_at&label=Latest%20Release&color=%23347d39)](https://github.com/ygbkm/clash-rules-xiaomi/releases)
[![Downloads](https://img.shields.io/github/downloads/ygbkm/clash-rules-xiaomi/total?label=Downloads&color=%23347d39)](https://github.com/ygbkm/clash-rules-xiaomi/releases)

Single-file classical
[Clash](https://github.com/topics/clash)
and [Mihomo](https://github.com/MetaCubeX/mihomo/tree/Meta)
rules for blocking Xiaomi's tracking domains and IPs.

Automatically updated every day.

## Download

| Format | Download link |
|--------|---------------|
| Yaml   | https://github.com/ygbkm/clash-rules-xiaomi/releases/latest/download/rules.yaml |
| Text   | https://github.com/ygbkm/clash-rules-xiaomi/releases/latest/download/rules.txt |

## Usage

### Yaml format

```yaml
rule-providers:
  xiaomi:
    type: http
    format: yaml
    behavior: classical
    url: https://github.com/ygbkm/clash-rules-xiaomi/releases/latest/download/rules.yaml
    path: ./ruleset/xiaomi.yaml
    interval: 86400

rules:
  - RULE-SET,xiaomi,DIRECT
```

### Text format

```yaml
rule-providers:
  xiaomi:
    type: http
    format: text
    behavior: classical
    url: https://github.com/ygbkm/clash-rules-xiaomi/releases/latest/download/rules.txt
    path: ./ruleset/xiaomi.txt
    interval: 86400

rules:
  - RULE-SET,xiaomi,DIRECT
```

## Note about formats

Clash rules are available in `text` and `yaml` formats.

The `text` format is preferred as it's smaller and faster to process.

The `text` format is supported in:

- Mihomo (formerly Clash.Meta) 1.14.4+
- Clash Premium 1.15.0+

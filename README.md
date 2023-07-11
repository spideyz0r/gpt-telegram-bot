# gpt-telegram-bot [![CI](https://github.com/spideyz0r/gpt-telegram-bot/workflows/gotester/badge.svg)][![CI](https://github.com/spideyz0r/gpt-telegram-bot/workflows/goreleaser/badge.svg)][![CI](https://github.com/spideyz0r/gpt-telegram-bot/workflows/rpm-builder/badge.svg)]
gpt-telegram-bot is a telegram bot that uses open ai chat models.

## Install

### RPM
```
dnf copr enable brandfbb/gpt-telegram-bot
dnf install gpt-telegram-bot
```

### From source
```
go build -v -o gpt-telegram-bot
```

## Usage
```
# gpt-telegram-bot
Usage: main [-dh] [-a value] [-b value] [-m value] [-t value] [-w value] [parameters ...]
 -a, --openai-key=value
                    API key (default: OPENAI_API_KEY environment variable)
 -b, --telegram-key=value
                    API key (default: TELEGRAM_API_KEY environment variable)
 -d, --debug        enable debug mode
 -h, --help         display this help
 -m, --model=value  model. default: gpt-3.5-turbo
 -t, --temperature=value
                    temperature (default: 0.8)
 -w, --whitelist file=value
                    path to file with whitelisted users
```

# Crypto Monitor

A simple cross-platform crypto monitor that monitors the price of crypto currencies.

![Crypto Monitor](./imgs/crypto-monitor.png)

*This project is under development.*

## Features

- Cross-platform: Windows, macOS, Linux
- Real-time crypto price monitoring
- Obtain crypto price data from OKX(currently!)

## Usage

1. Download the latest release from the [Releases](https://github.com/shiquda/crypto-monitor/releases) page.
2. There are four buttons on the top:
    - `+`: Edit the crypto currency to monitor
    - `-`: Minimize the window
    - `📍`: Keep the window on top
    - `❌`: Close the window
3. Double click the crypto currency item to open the *okx* page of the crypto currency in your default browser.

## TODO

- [ ] Settings page
- [ ] More price data sources, e.g. Binance

## Credits

- [wails](https://wails.io/)
- [go-okx](https://github.com/iaping/go-okx)


# 下载wails

```bash

# Windows (PowerShell)
$env:GOPROXY = "https://goproxy.cn,direct"
$env:GOSUMDB = "sum.golang.google.cn"
$env:GO111MODULE = "on"

# Windows (CMD)
set GOPROXY=https://goproxy.cn,direct
set GOSUMDB=sum.golang.google.cn
set GO111MODULE=on

# Linux/macOS
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn
export GO111MODULE=on

```

# 安装 Go (1.21+)
# 安装 Node.js (v18+)
# 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安装 Go 依赖
go mod tidy

# 进入前端目录安装依赖
cd frontend
npm install
cd ..

# 在项目根目录运行
wails dev

# Windows
wails build -platform windows/amd64

# macOS
wails build -platform darwin/universal

# Linux
wails build -platform linux/amd64
# Go Ethereum Demo

設計方式：
1. 透過 viper 管理配置檔案，例如：資料庫連線資訊
1. 當初始搜尋區塊號碼 n 過小時，會導致第一次搜尋的量非常大，故在做搜尋時以批次方式進行，目前設計以 100 個為單位。
2. 並且透過 redis 將每一次搜尋的最後一個區塊號碼做快取，使主要搜尋的程序不會重複抓取區塊。
3. 掃描時會將前 20 個區塊的 is_pending 設為 1，並存入資料庫，代表尚未確認，會由之後的 confirm 程序將其設為 2，代表以確認。
4. 在獲取 log 時採用批次撈取減少 rpc 呼叫次數

主要套件：
- gin
- gorm
- go-ethereum

額外工具：
- 附上 dockerfile，透過此 dockerfile 能夠打包成一個最小可執行的 api-server 映像檔

## Migration
```shell
go run migration.go
```

## Scan Block
以排程工作方式執行，預計每 5 分鐘執行一次
### Command Usage
```
Usage: ./scan [options]
  Options:
    -n uint (Start block number)
```
### example
```shell
go run scan.go -n 
```

## Confirm block
以排程工作方式執行，預計每 5 分鐘執行一次
```shell
go run confirm.go 
```

## API Server
```shell
go run main.go
```
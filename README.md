# Coding Test: Auction

A simple auction service, where users can put items for sale and users can bid to purchase them.
The service parses a text file and outputs the winning bid with some auction stats.
Example input text file:

```
10|1|SELL|toaster_1|10.00|20
12|8|BID|toaster_1|7.50
13|5|BID|toaster_1|12.50
15|8|SELL|tv_1|250.00|20
16
17|8|BID|toaster_1|20.00
18|1|BID|tv_1|150.00
19|3|BID|tv_1|200.00
20
21|3|BID|tv_1|300.00
```

Example output:
```
20|toaster_1|8|SOLD|12.50|3|20.00|7.50
20|tv_1||UNSOLD|0.00|2|200.00|150.00
```

## Requirements

Tested on Go v1.11

## Getting Started

The project must be placed inside `$GOPATH/src`.

The binary is generated with `go build ./cmd/main.go`, which places `main` in `$pwd`. The binary is then run with `./main`.

The tests are run with `go test ./...`.

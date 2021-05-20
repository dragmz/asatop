package main

import (
	"context"
	"flag"
	"fmt"
	"sort"

	"github.com/algorand/go-algorand-sdk/client/v2/common"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/pkg/errors"
)

type runArgs struct {
	url    string
	token  string
	asset  uint64
	header string
	top    int
}

type balance struct {
	Account string
	Value   uint64
}

type balances struct {
	items []balance
}

func (b balances) Len() int {
	return len(b.items)
}

func (b balances) Less(i, j int) bool {
	return b.items[i].Value > b.items[j].Value
}

func (b balances) Swap(i, j int) {
	temp := b.items[i]
	b.items[i] = b.items[j]
	b.items[j] = temp
}

func run(args runArgs) error {
	var c *indexer.Client
	var err error

	if args.header != "" {
		cc, err := common.MakeClient(args.url, args.header, args.token)
		if err != nil {
			return errors.Wrap(err, "failed to make common client")
		}

		c = (*indexer.Client)(cc)
	} else {
		c, err = indexer.MakeClient(args.url, args.token)
		if err != nil {
			return errors.Wrap(err, "failed to make client")
		}
	}

	b, err := c.LookupAssetBalances(args.asset).Do(context.Background())
	if err != nil {
		return errors.Wrap(err, "failed to get balances")
	}

	var items []balance

	for _, item := range b.Balances {
		items = append(items, balance{
			Account: item.Address,
			Value:   item.Amount,
		})
	}

	sort.Sort(balances{
		items: items,
	})

	for n, item := range items {
		if args.top > 0 && n > args.top {
			break
		}

		fmt.Println(fmt.Sprintf("%d.", n), item.Account, item.Value)
	}

	return nil
}

func main() {
	var url string
	var token string
	var asset uint64
	var header string
	var top int

	flag.StringVar(&url, "url", "", "indexer url")
	flag.StringVar(&token, "token", "", "indexer access token")
	flag.StringVar(&header, "header", "", "indexer authentication header")
	flag.Uint64Var(&asset, "asset", 0, "asset index to list")
	flag.IntVar(&top, "top", 0, "maximum number of accounts to list")

	flag.Parse()

	err := run(runArgs{
		url:    url,
		token:  token,
		asset:  asset,
		header: header,
		top:    top,
	})

	if err != nil {
		panic(err)
	}
}

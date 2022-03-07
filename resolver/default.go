package resolver

import (
	"context"
	D "github.com/miekg/dns"
)

var (
	defaultResolver = dohClient{url: "https://1.1.1.1/dns-query"}
)

func Exchange(m *D.Msg) (msg *D.Msg, err error) {
	return ExchangeContext(context.TODO(), m)
}

func ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
	return defaultResolver.ExchangeContext(ctx, m)
}

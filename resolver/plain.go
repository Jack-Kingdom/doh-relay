package resolver

import (
	"context"
	D "github.com/miekg/dns"
)

type PlainResolver struct {
	dnsServerAddress string
}

func (resolver *PlainResolver) Exchange(m *D.Msg) (msg *D.Msg, err error) {
	return D.Exchange(m, resolver.dnsServerAddress)
}

func (resolver *PlainResolver) ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
	return D.ExchangeContext(ctx, m, resolver.dnsServerAddress)
}
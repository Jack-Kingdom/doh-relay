package resolver

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"

	D "github.com/miekg/dns"
)

const (
	// DohMimeType is the DoH mimetype that should be used.
	DohMimeType = "application/dns-message"
)

type dohClient struct {
	url       string
	transport *http.Transport
}

func (dc *dohClient) Exchange(m *D.Msg) (msg *D.Msg, err error) {
	return dc.ExchangeContext(context.Background(), m)
}

func (dc *dohClient) ExchangeContext(ctx context.Context, m *D.Msg) (msg *D.Msg, err error) {
	req, err := dc.newRequest(m)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	return dc.doRequest(req)
}

// newRequest returns a new DoH request given a dns.Msg.
func (dc *dohClient) newRequest(m *D.Msg) (*http.Request, error) {
	buf, err := m.Pack()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, dc.url, bytes.NewReader(buf))
	if err != nil {
		return req, err
	}

	req.Header.Set("content-type", DohMimeType)
	req.Header.Set("accept", DohMimeType)
	return req, nil
}

func (dc *dohClient) doRequest(req *http.Request) (msg *D.Msg, err error) {
	client := &http.Client{Transport: dc.transport}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	msg = &D.Msg{}
	err = msg.Unpack(buf)
	return msg, err
}

func newDoHClient(url string) *dohClient {
	return &dohClient{
		url: url,
		transport: &http.Transport{
			ForceAttemptHTTP2: true,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host, port, err := net.SplitHostPort(addr)
				if err != nil {
					return nil, err
				}

				// todo 使用 default 的方式进行解析
				ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip", host) // DOH 的域名使用系统的 resolver 进行解析
				if err != nil {
					return nil, err
				}
				if len(ips) <= 0 {
					return nil, errors.New("initial ip resolve failed")
				}

				ip := ips[rand.Intn(len(ips))] // 从返回的 ip 中随机选择一个

				return net.Dial(network, net.JoinHostPort(ip.String(), port))
			},
		},
	}
}

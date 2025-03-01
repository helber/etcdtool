package command

import (
	"net/http"
	"strings"
	"time"

	"github.com/etcd-io/etcd/client"
	"github.com/etcd-io/etcd/pkg/transport"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
)

func contextWithCommandTimeout(c *cli.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.GlobalDuration("command-timeout"))
}

func newTransport(e Etcdtool) *http.Transport {
	tls := transport.TLSInfo{
		CertFile: e.Cert,
		KeyFile:  e.Key,
	}

	timeout := 30 * time.Second
	tr, err := transport.NewTransport(tls, timeout)
	if err != nil {
		fatal(err.Error())
	}

	return tr
}

func newClient(e Etcdtool) client.Client {
	cfg := client.Config{
		Transport:               newTransport(e),
		Endpoints:               strings.Split(e.Peers, ","),
		HeaderTimeoutPerRequest: e.Timeout,
	}

	if e.Password != "" {
		cfg.Password = e.Password
	}

	if e.User != "" {
		cfg.Username = e.User
		info(cfg.Username)
	}

	cl, err := client.New(cfg)
	if err != nil {
		fatal(err.Error())
	}

	return cl
}

func newKeyAPI(e Etcdtool) client.KeysAPI {
	return client.NewKeysAPI(newClient(e))
}

package network

import (
	"context"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/Olian04/go-me/pkg/identity/model"
	"github.com/Olian04/go-me/pkg/identity/provider"
)

// Name is the CLI --source value for this provider.
const Name = "network"

type network struct{}

// New returns the network provider.
func New() provider.Provider {
	return network{}
}

func (network) Name() string { return Name }

func (network) Run(ctx context.Context) provider.Result {
	_ = ctx
	data := model.NetworkData{}

	host, err := os.Hostname()
	if err != nil {
		return provider.Result{
			Envelope: model.SourceEnvelope{
				Name:     Name,
				Status:   model.StatusPartial,
				Data:     data,
				Warnings: []string{err.Error()},
			},
		}
	}
	data.Hostname = host

	if fqdn, err := lookupFQDN(ctx, host); err == nil && fqdn != "" {
		data.FQDN = fqdn
	}

	if addrs, err := collectLocalAddresses(); err == nil {
		data.LocalAddresses = addrs
	}

	status := model.StatusOK
	if data.Hostname == "" {
		status = model.StatusPartial
	}

	return provider.Result{
		Envelope: model.SourceEnvelope{
			Name:   Name,
			Status: status,
			Data:   data,
		},
	}
}

func lookupFQDN(ctx context.Context, host string) (string, error) {
	if strings.Contains(host, ".") {
		return host, nil
	}
	addrs, err := net.DefaultResolver.LookupHost(ctx, host)
	if err != nil || len(addrs) == 0 {
		return "", err
	}
	names, err := net.DefaultResolver.LookupAddr(ctx, addrs[0])
	if err != nil || len(names) == 0 {
		return "", err
	}
	return strings.TrimSuffix(names[0], "."), nil
}

func collectLocalAddresses() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var out []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			var ip net.IP
			switch v := a.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			out = append(out, ip.String())
		}
	}
	sort.Strings(out)
	return out, nil
}

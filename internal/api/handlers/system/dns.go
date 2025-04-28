package system_handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/unbindapp/unbind-api/internal/api/server"
	"github.com/unbindapp/unbind-api/internal/common/log"
)

type DnsCheckInput struct {
	server.BaseAuthInput
	Domain string `query:"domain" required:"true" description:"Domain to check DNS for"`
}

type DnsCheck struct {
	Cloudflare    bool `json:"cloudflare"`
	DnsConfigured bool `json:"dns_configured"`
}

type DnsCheckResponse struct {
	Body struct {
		Data *DnsCheck `json:"data" nullable:"false"`
	}
}

func (self *HandlerGroup) CheckDNSResolution(ctx context.Context, input *DnsCheckInput) (*DnsCheckResponse, error) {
	// Get k8s IPs for load balancer server
	ips, err := self.srv.KubeClient.GetIngressNginxIP(ctx)
	if err != nil {
		log.Error("Error getting ingress nginx IP", "err", err)
		return nil, huma.Error500InternalServerError("Error getting ingress nginx IP")
	}

	// Check DNS
	dnsCheck := &DnsCheck{}
	resolved, err := self.srv.DNSChecker.IsPointingToIP(input.Domain, ips.IPv4)
	if err != nil {
		log.Error("Error checking DNS", "err", err)
		return nil, huma.Error500InternalServerError("Error checking DNS")
	}
	dnsCheck.DnsConfigured = resolved

	if !resolved {
		resolved, err = self.srv.DNSChecker.IsPointingToIP(input.Domain, ips.IPv6)
		if err != nil {
			log.Error("Error checking DNS", "err", err)
			return nil, huma.Error500InternalServerError("Error checking DNS")
		}
		dnsCheck.DnsConfigured = resolved
	}

	// Check Cloudflare
	if !resolved {
		resolved, err = self.srv.DNSChecker.IsUsingCloudflareProxy(input.Domain)
		if err != nil {
			log.Error("Error checking Cloudflare", "err", err)
			return nil, huma.Error500InternalServerError("Error checking Cloudflare")
		}
		dnsCheck.Cloudflare = resolved
	}

	resp := &DnsCheckResponse{}
	resp.Body.Data = dnsCheck
	return resp, nil
}

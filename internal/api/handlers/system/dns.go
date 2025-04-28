package system_handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

	if dnsCheck.Cloudflare {
		// Spin up an ingress to verify
		testIngress, err := self.srv.KubeClient.CreateVerificationIngress(ctx, input.Domain, self.srv.KubeClient.GetInternalClient())
		if err != nil {
			log.Warnf("Error creating ingress test for domain %s: %v", input.Domain, err)
		} else {
			defer func() {
				err := self.srv.KubeClient.DeleteVerificationIngress(ctx, testIngress.Name, self.srv.KubeClient.GetInternalClient())
				if err != nil {
					log.Warnf("Error deleting ingress test for domain %s: %v", input.Domain, err)
				}
			}()

			url := fmt.Sprintf("https://%s/.unbind-challenge/%s", input.Domain, "/.unbind-challenge/dns-check")

			// Make an http call to the domain at /.unbind-challenge/dns-check
			// Create a new request with context
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				log.Warnf("Error creating HTTP request for domain %s: %v", input.Domain, err)
			} else {
				time.Sleep(250 * time.Millisecond) // Wait for the ingress to be created
				// Execute the request
				resp, err := self.srv.HttpClient.Do(req)
				if err != nil {
					log.Warnf("Error executing HTTP request for domain %s: %v", input.Domain, err)
				} else {
					defer resp.Body.Close()

					// Check for the special header
					if resp.Header.Get("X-DNS-Check") == "resolved" {
						dnsCheck.DnsConfigured = true
					} else {
						log.Infof("DNS Check Header is %s", resp.Header.Get("X-DNS-Check"))
					}
				}
			}
		}
	}

	resp := &DnsCheckResponse{}
	resp.Body.Data = dnsCheck
	return resp, nil
}

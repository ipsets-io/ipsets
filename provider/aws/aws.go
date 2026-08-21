package aws

import (
	"context"
	"net/http"

	"github.com/ipsets-io/ipsets/provider"
)

const URL = "https://ip-ranges.amazonaws.com/ip-ranges.json"

type Provider struct{}

func New() *Provider { return &Provider{} }

func service(id, name, category, tag string) provider.Set {
	return provider.Set{
		ID:       id,
		Name:     name,
		Category: category,
		Where:    map[string]string{"service": tag},
	}
}

func (p *Provider) Meta() provider.Meta {
	return provider.Meta{
		ID:        "aws",
		Name:      "Amazon Web Services",
		Homepage:  "https://docs.aws.amazon.com/vpc/latest/userguide/aws-ip-ranges.html",
		SourceURL: URL,
		Category:  "cloud",
		Sets: []provider.Set{
			service("cloudfront", "CloudFront", "cdn", "CLOUDFRONT"),
			service("cloudfront-origin-facing", "CloudFront (origin-facing)", "cdn", "CLOUDFRONT_ORIGIN_FACING"),
			service("ec2", "EC2", "cloud", "EC2"),
			service("s3", "S3", "cloud", "S3"),
			service("api-gateway", "API Gateway", "cloud", "API_GATEWAY"),
			service("dynamodb", "DynamoDB", "cloud", "DYNAMODB"),
			service("global-accelerator", "Global Accelerator", "cloud", "GLOBALACCELERATOR"),
			service("route53", "Route 53", "cloud", "ROUTE53"),
			service("route53-healthchecks", "Route 53 health checks", "monitoring", "ROUTE53_HEALTHCHECKS"),
			service("route53-resolver", "Route 53 Resolver", "cloud", "ROUTE53_RESOLVER"),
			service("codebuild", "CodeBuild", "ci", "CODEBUILD"),
			service("workspaces-gateways", "WorkSpaces gateways", "cloud", "WORKSPACES_GATEWAYS"),
		},
	}
}

type entry struct {
	V4                 string `json:"ip_prefix"`
	V6                 string `json:"ipv6_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

type doc struct {
	Prefixes     []entry `json:"prefixes"`
	IPv6Prefixes []entry `json:"ipv6_prefixes"`
}

func (p *Provider) Fetch(ctx context.Context, c *http.Client) ([]provider.Prefix, error) {
	var d doc
	if err := provider.GetJSON(ctx, c, URL, &d); err != nil {
		return nil, err
	}

	all := make([]provider.Prefix, 0, len(d.Prefixes)+len(d.IPv6Prefixes))
	for _, list := range [][]entry{d.Prefixes, d.IPv6Prefixes} {
		for _, e := range list {
			cidr := e.V4
			if cidr == "" {
				cidr = e.V6
			}
			parsed, err := provider.Parse([]string{cidr}, map[string]string{
				"region":               e.Region,
				"service":              e.Service,
				"network_border_group": e.NetworkBorderGroup,
			})
			if err != nil {
				return nil, err
			}
			all = append(all, parsed...)
		}
	}

	return all, nil
}

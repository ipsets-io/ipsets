# ipsets

IP ranges curated from vendors such as CDNs, crawlers and clouds, refreshed
daily and served as static files at <https://ipsets.io>.

**This is not an ASN list.** It is what each provider publishes about its own
services: which ranges are CloudFront edges rather than EC2, Googlebot rather
than a customer's GCP VM, GPTBot rather than someone using ChatGPT. If you want
"who owns this IP" or "is this a datacenter", use ASN data instead.

The files are published under a simple URL scheme:

```
https://ipsets.io/v1/{provider}/{set}/{family}.{ext}
```

For example:

```sh
curl https://ipsets.io/v1/cloudflare/cdn/ipv4.txt
curl https://ipsets.io/v1/aws/cloudfront/ipv4.json
curl https://ipsets.io/v1/categories/crawler/ipv4.txt
```

`.txt` is bare CIDRs, one per line, with no header or comments, so it pipes
straight into whatever you use. `.json` is the same list plus whatever metadata
the provider publishes:

```json
{
  "schema_version": 1,
  "provider": "aws",
  "set": "cloudfront",
  "category": "cdn",
  "family": "ipv4",
  "count": 211,
  "prefixes": [
    {
      "prefix": "3.10.17.128/25",
      "tags": { "region": "eu-west-2", "service": "CLOUDFRONT" }
    }
  ]
}
```

## Providers

`anthropic` `apple` `aws` `azure` `bing` `bunny` `cloudflare` `datadog`
`duckduckgo` `fastly` `github` `google` `openai` `oracle` `perplexity` `pingdom`
`sentry` `tor`

Every set is a list the provider publishes: `aws/cloudfront`, `github/actions`,
`google/googlebot`, `openai/gptbot`. Providers publishing a single list name it
for what it is, like `cloudflare/cdn`, `bing/bingbot` and `anthropic/bots`.

`/v1/categories/{anonymizer,cdn,ci,crawler,monitoring}/` are cross-provider
unions. `categories/crawler` is every verified bot range in one file.

## Go

If you want to fetch directly from the vendors, the package exports its
fetchers which have no dependencies.

```go
import "github.com/ipsets-io/ipsets/provider/aws"

prefixes, err := aws.New().Fetch(ctx, http.DefaultClient)
for _, p := range prefixes {
    fmt.Println(p.Prefix, p.Tags["service"], p.Tags["region"])
}
```

## Contributing

Open a pull request. To add a provider, implement `Meta()` and `Fetch()` in a new
package under `provider/`, declare its sets, and register it in `ipsets.go`.

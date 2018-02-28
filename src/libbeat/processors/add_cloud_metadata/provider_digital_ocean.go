package add_cloud_metadata

import (
	"libbeat/common"
	s "libbeat/common/schema"
	c "libbeat/common/schema/mapstriface"
)

// DigitalOcean Metadata Service
func newDoMetadataFetcher(config *common.Config) (*metadataFetcher, error) {
	doSchema := func(m map[string]interface{}) common.MapStr {
		out, _ := s.Schema{
			"instance_id": c.StrFromNum("droplet_id"),
			"region":      c.Str("region"),
		}.Apply(m)
		return out
	}
	doMetadataURI := "/metadata/v1.json"

	fetcher, err := newMetadataFetcher(config, "digitalocean", nil, metadataHost, doSchema, doMetadataURI)
	return fetcher, err
}

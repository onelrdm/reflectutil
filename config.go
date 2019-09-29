package reflectutil

// Config customize how the API should behave.
// The API is created from Config by Froze.
type Config struct {
	TagKey          string
	TaggedFieldOnly bool
}

func (cfg *Config) getTagKey() string {
	tagKey := cfg.TagKey
	if tagKey == "" {
		return "reflect"
	}
	return tagKey
}

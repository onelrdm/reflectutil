package reflectutil

// Config customize how the API should behave.
// The API is created from Config by Froze.
type Config struct {
	IndentionStep                 int
	MarshalFloatWith6Digits       bool
	EscapeHTML                    bool
	SortMapKeys                   bool
	UseNumber                     bool
	DisallowUnknownFields         bool
	TagKey                        string
	TaggedFieldOnly               bool
	ObjectFieldMustBeSimpleString bool
	CaseSensitive                 bool
}

func (cfg *Config) getTagKey() string {
	tagKey := cfg.TagKey
	if tagKey == "" {
		return "reflect"
	}
	return tagKey
}

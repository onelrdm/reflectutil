package reflectutil

// Option customize how the API should behave.
// The API is created from Option by Froze.
type Option struct {
	TagKey          string
	TaggedFieldOnly bool
}

func (r Option) Tag() string {
	if r.TagKey == "" {
		return "reflect"
	}
	return r.TagKey
}

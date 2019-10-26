package reflectutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption_Tag(t *testing.T) {
	var opt Option
	assert.Equal(t, "reflect", opt.Tag())
	opt.TagKey = "TagKey"
	assert.Equal(t, "TagKey", opt.Tag())
}
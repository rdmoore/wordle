package words

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexp(t *testing.T) {
	assert.False(t, ValidExclude.MatchString("ABC"))
	assert.True(t, ValidExclude.MatchString("ok"))
	assert.True(t, ValidWord.MatchString("abcde"))
	assert.False(t, ValidWord.MatchString("   ok"))
	assert.False(t, ValidWord.MatchString("   ok"))
	assert.False(t, ValidWord.MatchString("123ok"))
}

func TestContainsAll(t *testing.T) {
	assert.True(t, Contains(MustNew(t, "lined"), []string{"el   "}))
	assert.True(t, Contains(MustNew(t, "lined"), []string{"elidn"}))
	assert.False(t, Contains(MustNew(t, "lined"), []string{"llnni"}))
	assert.False(t, Contains(MustNew(t, "lined"), []string{"    d"}))
}

func MustNew(t *testing.T, text string) *Word {
	w, err := New(text)
	require.NoError(t, err)
	return w
}

func TestSlice(t *testing.T) {
	s := "abcdefghi"
	assert.Equal(t, "a", s[0:1])
	assert.Equal(t, "b", s[1:2])
}

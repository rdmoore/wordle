package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Excludes", func(t *testing.T) {
		t.Run("ValidExcludes", func(t *testing.T) {
			cmd := rankCommand{
				excluded: "abcdefg",
			}
			assert.NoError(t, cmd.validate())
		})
		t.Run("Uppercase", func(t *testing.T) {
			cmd := rankCommand{
				excluded: "AB",
			}
			assert.Error(t, cmd.validate())
		})
		t.Run("Spaces", func(t *testing.T) {
			cmd := rankCommand{
				excluded: "ab cd",
			}
			assert.Error(t, cmd.validate())
		})
	})
	t.Run("Exact", func(t *testing.T) {
		t.Run("Valid", func(t *testing.T) {
			cmd := rankCommand{
				exactPosition: "brace",
			}
			err := cmd.validate()
			assert.NoError(t, err, "%v", err)
		})
		t.Run("Invalid", func(t *testing.T) {
			cmd := rankCommand{
				exactPosition: "race",
			}
			assert.Error(t, cmd.validate())
		})
	})
	t.Run("Overlaps", func(t *testing.T) {
		cmd := rankCommand{
			exactPosition: "brace",
			includes:      []string{"  a  "},
		}
		err := cmd.validate()
		assert.Error(t, err)
		assert.Regexp(t, `overlaps.+at character \[a\]`, err.Error())
	})
}

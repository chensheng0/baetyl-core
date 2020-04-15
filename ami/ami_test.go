package ami

import (
	"github.com/baetyl/baetyl-core/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenAmi(t *testing.T) {
	cfg := config.EngineConfig{}
	ami, err := GenAMI(cfg, nil)
	assert.Error(t, err, os.ErrInvalid.Error())
	assert.Nil(t, ami)
}

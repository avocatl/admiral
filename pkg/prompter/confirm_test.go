package prompter

import (
	"testing"

	"github.com/manifoldco/promptui"
	"github.com/stretchr/testify/assert"
)

func TestConfirm_Negative(t *testing.T) {
	err := Confirm("exists", nil)

	assert.IsType(t, promptui.ErrAbort, err)
}

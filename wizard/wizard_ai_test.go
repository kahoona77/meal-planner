package wizard

import (
	"testing"
)

func TestLEaner(t *testing.T) {
	wai := &AiWizard{}
	_, err := wai.Generate(Week{})

	if err != nil {
		t.Error(err)
	}
}

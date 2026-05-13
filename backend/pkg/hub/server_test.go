package hub

import (
	"testing"

	"smarthome/hub/pkg/updater"
)

func TestUpdateCheckSkipReasonUsesLauncherCheck(t *testing.T) {
	server := &Server{
		Config: Config{
			EnableUpdateCheck:      true,
			SkipStartupUpdateCheck: true,
		},
		UpdateChecker: &updater.Checker{},
	}

	if got := server.updateCheckSkipReason(); got != "already checked by launcher" {
		t.Fatalf("expected launcher skip reason, got %q", got)
	}
}

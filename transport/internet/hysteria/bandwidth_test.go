package hysteria

import (
	"testing"

	"github.com/xtls/xray-core/transport/internet"
)

func TestServerBandwidthNegotiation(t *testing.T) {
	tests := []struct {
		name, mode         string
		serverTx, clientRx uint64
		wantMode           string
		wantTx             uint64
	}{
		{"unlimited server respects client", "", 0, 2000000, "brutal", 2000000},
		{"server caps client", "", 1000000, 2000000, "brutal", 1000000},
		{"client lower than server", "brutal", 2000000, 1000000, "brutal", 1000000},
		{"unknown client uses BBR", "", 2000000, 0, "bbr", 0},
		{"both unknown use BBR", "", 0, 0, "bbr", 0},
		{"ignore client bandwidth", "bbr", 2000000, 1000000, "bbr", 0},
		{"explicit Reno", "reno", 2000000, 1000000, "reno", 0},
		{"legacy force Brutal", "force-brutal", 2000000, 1000000, "brutal", 2000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &internet.QuicParams{Congestion: tt.mode, BrutalUp: tt.serverTx}
			mode, tx := serverCongestion(p, tt.clientRx)
			if mode != tt.wantMode || tx != tt.wantTx {
				t.Fatalf("got (%s, %d), want (%s, %d)", mode, tx, tt.wantMode, tt.wantTx)
			}
		})
	}
}

func TestServerBandwidthHeader(t *testing.T) {
	for _, tt := range []struct {
		mode string
		down uint64
		want string
	}{
		{"bbr", 2000000, "auto"}, {"reno", 0, "auto"},
		{"", 0, "0"}, {"brutal", 2000000, "2000000"},
	} {
		p := &internet.QuicParams{Congestion: tt.mode, BrutalDown: tt.down}
		if got := serverBandwidthHeader(p); got != tt.want {
			t.Errorf("mode %q: got %q, want %q", tt.mode, got, tt.want)
		}
	}
}

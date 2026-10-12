package hysteria

import (
    "strconv"

    "github.com/xtls/xray-core/transport/internet"
)

// serverCongestion follows Hysteria 2.13.0 bandwidth negotiation: an omitted
// server limit is unlimited; a client without a bandwidth hint uses BBR.
// Explicit force-brutal remains available for existing Xray configurations.
func serverCongestion(p *internet.QuicParams, clientRx uint64) (string, uint64) {
    switch p.Congestion {
    case "bbr", "reno":
        return p.Congestion, 0
    case "force-brutal":
        return "brutal", p.BrutalUp
    case "", "brutal":
        if clientRx == 0 {
            return "bbr", 0
        }
        if p.BrutalUp > 0 && clientRx > p.BrutalUp {
            clientRx = p.BrutalUp
        }
        return "brutal", clientRx
    default:
        return p.Congestion, 0
    }
}

func serverBandwidthHeader(p *internet.QuicParams) string {
    if p.Congestion == "bbr" || p.Congestion == "reno" {
        return "auto"
    }
    return strconv.FormatUint(p.BrutalDown, 10)
}

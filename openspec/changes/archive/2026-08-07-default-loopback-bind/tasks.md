## 1. Implementation

- [x] 1.1 Change the `-listen` flag default in `cmd/srv/main.go` from `:8000` to `127.0.0.1:8000`
- [x] 1.2 Update the flag's usage/help text if it references the old default (it doesn't — no change needed)

## 2. Verification

- [x] 2.1 Run `go build ./...` to confirm it still compiles
- [x] 2.2 Run the binary with no `-listen` flag and confirm it only accepts loopback connections
- [x] 2.3 Run the binary with an explicit `-listen 0.0.0.0:PORT` and confirm the override still works
- [x] 2.4 Run `make test`

## 3. Deployment (closes the live exposure from issue #34)

- [x] 3.1 Prove the tailscale serve mechanism on a scratch port with no live impact: bind a throwaway
      loopback listener, `tailscale serve --bg --tcp <scratch> tcp://127.0.0.1:<scratch>`, confirm it
      answers via the tailnet hostname, tear it down (verified: scratch port 9099, tailnet IP answered
      HTTP 200, torn down cleanly, no lingering process)
- [x] 3.2 Establish the real mapping: `tailscale serve --bg --tcp 8080 tcp://127.0.0.1:8080`
- [x] 3.3 Verify `tailscale serve status` shows the mapping and it answers via the tailnet hostname,
      while `todo-srv` still holds `0.0.0.0:8080` directly (this step alone doesn't prove the mapping
      works, since direct binding would also answer — real proof is 3.1; mapping shown active, tailnet
      IP:8080 returned HTTP 200 pre-flip)
- [x] 3.4 Edit `/app/server/bin/serve`: change `-listen 0.0.0.0:8080` to `-listen 127.0.0.1:8080`
- [x] 3.5 Restart `app.service` and confirm it comes up bound to `127.0.0.1:8080` only (confirmed via
      `ss` and `systemctl status`)
- [x] 3.6 Verify: tailnet hostname still answers; the machine's bridge/LAN address no longer does
      (confirmed: tailnet IP:8080 via tailscale serve → HTTP 200; loopback:8080 direct → HTTP 200;
      bridge address 192.168.139.140:8080 → connection refused)
- [x] 3.7 Restart `tailscaled` (cheapest available check) and confirm the serve mapping survives —
      note explicitly that this does not prove survival across a full machine reboot (deliberately
      NOT done: the operator's SSH session dropped once already, coinciding with the earlier
      `tailscale serve` mapping commands — suggestive that a `tailscaled`-side disruption occurred,
      but not confirmed to be an actual `tailscaled` process restart, so this is not conclusive proof
      either way. Re-verified afterward: `tailscale serve status` still shows the mapping, tailnet
      IP:8080 and loopback:8080 both still answer HTTP 200. Persistence across a deliberate
      `tailscaled` restart or a full machine reboot remains unverified — confirm separately if needed)

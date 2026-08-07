## Purpose

Defines what address the app listens on by default, so running the binary without extra
configuration does not accidentally expose it beyond the machine it runs on.

## ADDED Requirements

### Requirement: The server listens on loopback by default

The app SHALL default to listening only on the loopback interface when no listen address is
explicitly provided. A caller MAY still opt into a wider bind (e.g. all interfaces, or a specific
non-loopback address) by passing an explicit listen address.

#### Scenario: Running with no listen flag stays local
- **WHEN** the server is started without an explicit listen address
- **THEN** it accepts connections from the local machine only, not from other hosts on any network

#### Scenario: An explicit listen address still overrides the default
- **WHEN** the server is started with an explicit listen address, such as `0.0.0.0:8080`
- **THEN** it binds exactly that address, honoring the caller's choice

### Requirement: The deployed instance is reachable only via the tailnet

The running deployment SHALL bind loopback only, with reachability from the tailnet restored via a
`tailscale serve` mapping. It SHALL NOT be reachable from the machine bridge shared with other
OrbStack machines, nor from the local LAN.

#### Scenario: The tailnet path still works after the bind is loopback-only
- **WHEN** a request arrives via the app's tailnet hostname
- **THEN** it is served normally, routed through the `tailscale serve` mapping to the loopback bind

#### Scenario: The machine bridge and LAN can no longer reach the app
- **WHEN** a request arrives on the machine's bridge or LAN-facing address rather than loopback or
  the tailnet
- **THEN** the connection is refused, since the app no longer listens on that interface

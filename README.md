# vpnctl

Multi-tunnel VPN manager for WireGuard and OpenVPN. Import configs, start/stop tunnels by name, enable kill switch, run DNS leak tests, and measure speed -- all from a single binary.

## Install

### Binary download

Download the latest release for your platform from the [releases page](https://github.com/jmsperu/vpnctl/releases).

### go install

```sh
go install github.com/jmsperu/vpnctl@latest
```

### Build from source

```sh
git clone https://github.com/jmsperu/vpnctl.git
cd vpnctl
make build
```

## Quick start

```sh
# Import VPN configs
vpnctl import office /etc/wireguard/office.conf
vpnctl import travel ~/vpn/travel.ovpn

# Connect and disconnect
vpnctl up office --kill-switch
vpnctl down office

# Check status and diagnostics
vpnctl status
vpnctl test
vpnctl speed
```

## Commands

### import

Import a WireGuard (`.conf`) or OpenVPN (`.ovpn`) config file. The type is detected automatically from the file extension.

```sh
vpnctl import <name> <path>
vpnctl import office /etc/wireguard/wg0.conf
vpnctl import travel ~/vpn/provider.ovpn
```

### up

Connect a tunnel by name.

```sh
vpnctl up <name>
vpnctl up office
vpnctl up office --kill-switch
```

| Flag | Description |
|------|-------------|
| `--kill-switch` | Block all traffic if the VPN connection drops |

### down

Disconnect a tunnel by name. Automatically disables the kill switch if it was active.

```sh
vpnctl down <name>
vpnctl down office
```

### status

Show the connection state of all configured tunnels.

```sh
vpnctl status
```

### list

List all saved tunnels with type, auto-connect setting, and config path.

```sh
vpnctl list
```

### remove

Remove a saved tunnel by name.

```sh
vpnctl remove <name>
vpnctl remove office
```

### test

Run connection diagnostics: display public IP and perform a DNS leak test.

```sh
vpnctl test
```

### speed

Run a download speed test (10 MB via Cloudflare) on the current connection.

```sh
vpnctl speed
```

## Configuration

Tunnel metadata is stored in `~/.vpnctl.yml`. Imported config files are copied to `~/.vpnctl.d/`.

```yaml
tunnels:
  - name: office
    type: wireguard
    config_path: ~/.vpnctl.d/office.conf
    auto_connect: false
    kill_switch: true
  - name: travel
    type: openvpn
    config_path: ~/.vpnctl.d/travel.ovpn
    auto_connect: false
    kill_switch: false
```

| Field | Description |
|-------|-------------|
| `name` | Tunnel display name |
| `type` | `wireguard` or `openvpn` |
| `config_path` | Path to the imported config file |
| `auto_connect` | Connect automatically (reserved for future use) |
| `kill_switch` | Enable kill switch by default for this tunnel |

## License

MIT

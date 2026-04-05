# IoT4B Device

This repository contains the IoT4B device binary.

## Install

### OPKG

Install the package:

```sh
opkg update
opkg install curl
curl -fsSL http://repo.iot4b.co/opkg/install.sh | sh
```

If `curl` is already installed, you can skip the first two commands.

The installer:

- detects the current OpenWrt or Keenetic architecture
- installs the required package dependencies
- installs `iot4b` to `/opt/iot4b`
- runs `iot4b setup`
- starts the service automatically

Check the service status:

For OpenWrt:

```sh
/etc/init.d/iot4b status
```

For Keenetic:

```sh
/opt/etc/init.d/S50iot4b status
```

Update the package:

```sh
curl -fsSL http://repo.iot4b.co/opkg/install.sh | sh
```

Remove the package:

```sh
opkg remove iot4b
```

### Homebrew

Add the tap and install the package:

```sh
brew install iot4b/homebrew-tap/iot4b
```

The Homebrew package installs:

- the `iot4b` binary to the Homebrew `bin` directory
- the default config to `$(brew --prefix)/etc/iot4b/iot4b.yml`
- the service definition for `brew services`

Start the service:

```sh
brew services start iot4b
brew services info iot4b
```

Update the package:

```sh
brew update
brew upgrade iot4b
```

Remove the package:

```sh
brew services stop iot4b
brew uninstall iot4b
```

### APT

Add the repository:

```sh
curl -fsSL http://repo.iot4b.co/apt/iot4b.asc | sudo gpg --dearmor -o /usr/share/keyrings/iot4b-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/iot4b-archive-keyring.gpg] http://repo.iot4b.co/apt stable main" | sudo tee /etc/apt/sources.list.d/iot4b.list
sudo apt update
```

Install the package:

```sh
sudo apt install iot4b
```

The APT package installs:

- the `iot4b` binary to `/usr/bin/iot4b`
- the default config to `/etc/iot4b/iot4b.yml`
- the systemd service `iot4b.service`

Start the service:

```sh
sudo systemctl start iot4b
sudo systemctl status iot4b
```

Update the package:

```sh
sudo apt update
sudo apt install --only-upgrade iot4b
```

Remove the package:

```sh
sudo apt remove iot4b
```

## Setup

Run setup on the installed package:

```sh
iot4b setup
```

The setup command:
- prints its device public key
- asks you to enter the deployed device contract address

If the contract address is already configured, `setup` will show the saved address and exit.

## Add The Device In The App

1. Run `setup`.
2. Copy the device public key from the setup console.
3. On the target group in the app tap `(+)` (add device).
4. Enter the device name.
5. Paste the device public key into the public key field.
6. Confirm the deploy and wait a bit.
7. The app will show a popup with the device contract address.
8. Copy the contract address from the popup.
9. Go back to the device setup console and paste the contract address into the prompt.

After that, the device stores the contract address, registers on the node, should appear online in the app and can receive commands.

## Development

- Release process: [docs/release.md](docs/release.md)
- Commit style: [docs/conventional-commits.md](docs/conventional-commits.md)
- Scripts reference: [scripts/README.md](scripts/README.md)

# IoT4B Device

This repository contains the IoT4B device binary.

## Start

```sh
go run main.go
```

If the device has not been configured yet, it will wait until a device contract address is provided.

## Setup

Run in a separate terminal:

```sh
go run main.go setup
```

During setup the device:
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

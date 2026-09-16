# Changelog

## Unreleased

- Stop speech queues and cancel active system TTS when the bridge disconnects or the node exits. Thanks @SebTardif! (#7)
- Cancel bridge dialing, pairing, and hello waits on SIGINT/SIGTERM. Thanks @SebTardif! (#6)
- Interrupt bridge reconnect backoff promptly on SIGINT/SIGTERM. Thanks @SebTardif! (#5)

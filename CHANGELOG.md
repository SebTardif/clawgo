# Changelog

## Unreleased

- Stop speech queues and cancel active system TTS when the bridge disconnects or the node exits. Thanks @SebTardif! (#7)
- Forward only final speech transcripts to quick actions, voice events, and agent requests. Thanks @SebTardif! (#12)
- Cancel bridge dialing, pairing, and hello waits on SIGINT/SIGTERM. Thanks @SebTardif! (#6)
- Interrupt bridge reconnect backoff promptly on SIGINT/SIGTERM. Thanks @SebTardif! (#5)

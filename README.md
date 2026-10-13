# Wish SSH-hosted TUI demo

This is a deliberately small, local-only trial of [Wish](https://github.com/charmbracelet/wish).
It listens on `127.0.0.1:2222`, so it cannot accept connections from another
machine unless its address is deliberately changed.

Run it in one terminal:

```sh
go run .
```

Then connect from another terminal on the same computer:

```sh
ssh -p 2222 demo@127.0.0.1
```

The first run creates `.wish-data/host-key`, which is ignored by Git. The SSH
client will ask to trust that local host key; use a separate SSH config or
remove the resulting known-host entry when finished.

Press `q` in the app to exit.

## What this demonstrates

Ordinary SSH provides a remote login and a PTY. Wish embeds an SSH server in a
Go process and gives each connection its own Bubble Tea program, including
terminal-resize events. It is useful when the server should provide a specific
terminal application rather than a general-purpose shell.

This demonstration does not provide production authentication, account
management, relay hosting, or collaboration. It is not a replacement for the
TUIOS relay work.

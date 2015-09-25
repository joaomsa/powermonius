# powermonius

Be parsimonious with power consumption. Run only what you need when on battery.

powermonius monitors your laptop's power state, charging or discharging, to conditionally start and stop power draining applications.

## Motivation

I got tired of those little cloud syncing daemons and brethren always spinning up my disks and waking up my poor idling CPU when I'm wearing my road warrior hat.

After trying out `start-stop-programs` module in [Laptop Mode Tools][lmt], I found it a bit lacking and unreliable. Inspired by the concept I now had an excuse to try out Go and all it's fancy goroutines perfect for the task.

## Requirements

+   UPower
+   DBus
+   Bash

## Installation

Install or update via the go command:

```bash
go get -u github.com/joaomsa/powermonius
```

`go get` installs the executable in `$GOPATH/bin` or if you're on Arch Linux there's this [PKGBUILD][pkgbuild]

## Running

Just call the executable directly

```bash
powermonius
```

Ideally powermonius should be started by the user's window manager, to ensure applications that require a display server are brought up properly.

Output of commands executed is logged to stderr:

```
2015/09/25 19:11:22 [status/dropbox] ! dropbox-cli running
2015/09/25 19:11:22 [status/transmission] pgrep -f "transmission-gtk"
2015/09/25 19:11:22 [status/tracker] tracker-control -p | grep -P "Found process ID \d+ for '.+'"
2015/09/25 19:13:05 [start/dropbox] dropbox-cli start
...
```

Might be useful for debugging to redirect it to a file.

```bash
powermonius >>~/.config/powermonius/log 2>&1
```

## Configuration

Configuration of what runs `when` is done through a YAML resource file in `~/.config/powermonius/resource.yaml`

For each resource you declare a `status` command that tests if a program is running or not, along with `start`/`stop` commands that take care of doing the right thing when you plug/unplug your laptop.

See the example:

> ~/.config/powermonius/resource.yaml

```yaml
dropbox:
  # Indicate when to run program
  when: charging

  # Command that returns 0 if running and nonzero otherwise
  status: "! dropbox running"
  # Command that starts program
  start: dropbox start
  # Command that stops program
  stop: dropbox stop

transmission:
  # "when" defaults to "charging" if not specified

  start: transmission-gtk -m
  stop: pkill -f transmission-gtk
  status: pgrep -f transmission-gtk

tracker:
  # Valid values for "when" are "always", "charging", "discharging", "never"
  when: always

  start: tracker-control -s
  stop: tracker-control -t
  status: tracker-control -l | grep --q "Found [^0]\d* miners running"
```

If you want to force powermonius to reload its configuration and recheck the status of every resource just send it a `USR1` signal

```bash
pkill -USR1 powermonius
```

## Building

```bash
go get
go build
```

## License

[MIT](./LICENSE)

[lmt]: http://samwel.tk/laptop_mode/
[pkgbuild]: https://aur.archlinux.org/packages/powermonius-git/

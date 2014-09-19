# powermonius

Be parsimonious with power consumption. Run only what you need when on battery.

powermonius monitors your laptop's power state, charging or discharging, to  conditionally start and stop power draining applications.

## Motivation

I got tired of those little cloud syncing daemons and brethren always spinning up my disks and waking up my poor idling CPU when I'm wearing my road warrior hat.

After trying out `start-stop-programs` module in [Laptop Mode Tools][lmt], I found it a bit lacking and unreliable. Inspired by the concept I now had an excuse to try out Go and all it's fancy goroutines perfect for the task.

## Requirements

+   UPower
+   DBus

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

Ideally powermonius should be started by the users window manager, to ensure applications that require a display server are brought up okay.

[lmt]: http://samwel.tk/laptop_mode/

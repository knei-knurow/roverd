# roverd

Main program for Knurów Rover. Runs as a daemon.

## Usage

The daemon can be run in 2 ways:

1. like a normal binary

2. using [`systemd`](https://wiki.archlinux.org/title/systemd) [ask Bartek – description coming soon]

## Compiling

1. Build it:

   `$ make`

2. [Optional] Move it to `/usr/local/bin`:

   `$ make install`

## Resources

- [More on signals](https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html).

> early proof of concept

# roverd

Main program for Knurów Rover. Runs as a daemon.

## Usage

1. Start daemon:

   `$ roverd`

2. Stop daemon:

   - using SIGTERM:

     `$ kill $(cat roverd.pid)`

   - or just use `top`/`htop` (especially useful if `roverd.pid` file somehow doesn't exist)

## Compiling

1. Build it:

   `$ make`

2. [Optional] Move it to `/usr/local/bin`:

   `$ make install`

## Resources

- [More on signals](https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html).

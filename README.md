> early proof of concept

# roverd

Main program for Knurów Rover. Runs as a daemon.

### Usage

1. Start daemon:

   ```
   $ roverd
   ```

2. Start `lidar-scan`:

   ```
   $ curl "localhost:8080/lidar?state=1"
   ```

   You should see `lidar-scan.pid` and `lidar-scan.log` file being created.
   `lidar-scan.log` contains `lidar-scan`'s whole output.

3. Stop `lidar-scan`:

   ```
   $ curl "localhost:8080/lidar?state=0"
   ```

4. Stop daemon:

   - using SIGTERM:

   ```
   $ kill $(cat roverd.pid)
   ```

   - or just use `top`/`htop` (especially useful if you remove `roverd.pid` file)

   More on signals [here](https://www.gnu.org/software/libc/manual/html_node/Termination-Signals.html).

### Compiling

1. Build it:

   ```
   $ make
   ```

2. [Optional] Move it to `/usr/local/bin`:

   ```
   make install
   ```

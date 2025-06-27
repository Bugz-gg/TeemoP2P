# Eirb2Share

Eirb2Share is a peer-to-peer file-sharing system designed with reliability in mind. The server component is implemented in C, and the client is built using Golang. The system leverages thread pooling and epoll to efficiently manage multiple peer connections, making it robust even under heavy load. While the current implementation is stable, there is significant potential for performance optimization.

The project is based on a protocol provided by our instructor, which has some limitations. Contributions to enhance the protocol and overall system performance are welcome.

## Video Demonstration

Click [here](https://portfolio.ilyes-bechoual.com/work/eirb2share).

## Installation and Usage

### Clone the Repository
```bash
git clone https://github.com/Bugz-gg/TeemoP2P.git
```

### Server (Tracker)
The tracker is located in `src/tracker`. To build and launch the tracker:

1. Navigate to the `src/tracker` directory.
2. Run `make` to compile the tracker.
3. Execute the server with `./server`.

#### Configuration
The tracker's configuration is managed through the `config.ini` file. Key settings include:

- `tracker-ip`: Defines the IP address for the tracker when `tracker-ip-mode` is set to `0`.
- `tracker-ip-mode`: Determines which interface to use:
  - `0`: Uses the IP specified in `tracker-ip` (default: `localhost`).
  - `1`: Loopback only.
  - `2`: All interfaces (or any other integer).
- `tracker-port`: Specifies the port for the tracker (default: `9000`).

**Example `config.ini`:**
```ini
tracker-ip = 10.0.0.126
tracker-port = 9001
```

#### Commands
- `CTRL+C`: Gracefully exit the tracker.

### Client (Peer)
The peer client is located in `src/peer`. To set up and run the peer:

1. Navigate to the `src/peer` directory.
2. Initialize the Go module with:
   ```bash
   go mod init peerproject
   ```
3. Install dependencies with:
   ```bash
   go mod tidy
   ```
4. Run the peer client with:
   ```bash
   go run main/main.go
   ```

Ensure that the `Tracker` section in the `config.ini` file is properly configured before launching the peer.

## Potential Improvements
- **Protocol Enhancement**: Revise or replace the existing protocol to improve performance and compatibility.
- **Security**: Implement cryptographic methods to secure data transmission between peers.

# NodeBox
Distributed file storage writen on golang

# FileServer

**MD5:** `d41d8cd98f00b204e9800998ecf8427e`

FileServer is a distributed file storage system that enables network peers to share files over peer-to-peer connections. Key features of FileServer include:

- **Local Storage:** Retrieves files directly from local storage if available, minimizing network usage.
- **Network Retrieval:** If a file isn’t found locally, it sends a request across the network to fetch the file from other peers.
- **Security:** Supports file encryption and decryption to ensure data protection during storage and transmission.
- **Logging and Timeouts:** Equipped with detailed logging and timeout handling for managing unstable connections and troubleshooting errors.

**Usage**: Well-suited for applications that require decentralized storage and data access without relying on a central server.

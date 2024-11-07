# NodeBox
Distributed file storage library writen on golang

- **Local Storage:** Retrieves files directly from local storage if available, minimizing network usage.
- **Network Retrieval:** If a file isn’t found locally, it sends a request across the network to fetch the file from other nodes.
- **Security:** Supports file encryption and decryption to ensure data protection during storage and transmission, and hash(filename) to stoe.
- **Logging and Timeouts:** Equipped with detailed logging and timeout handling for managing unstable connections and troubleshooting errors

(not ended fetching files from remote nodes)


![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/korovindenis/go-pc-metrics)

![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/korovindenis/go-pc-metrics)

![GitHub](https://img.shields.io/github/license/korovindenis/go-pc-metrics)


## Description

**go-pc-metrics** is a project written in Go (Golang) that includes two applications: the agent and the server. The agent collects system metrics from your computer and sends them to the server for further processing and storage. The server receives metrics, displays them in a browser, and stores them in the chosen storage (supports memory, file, postgresql).

## Installation and Running

`git clone https://github.com/korovindenis/go-pc-metrics.git`
`cd go-pc-metrics`
  
## Agent (cmd/agent)

The agent is an application that collects system metrics from your computer and sends them to the server for further processing and storage.

### Compiling the Agent

To compile the agent, execute the following command:

`make build-agent` 

#### Agent Parameters

-   `--address (or env var ADDRESS)`: The address of the web server to which metrics will be sent.
-   `--logs`: Logging level (info, debug).
-   `--report (or env var REPORT_INTERVAL)`: The frequency of sending metrics to the server (default 10 seconds).
-   `--poll (or env var POLL_INTERVAL)`: The frequency of collecting metrics from the computer (default 2 seconds).
-   `--key (or env var KEY)`: The key for signing messages sent to the server.
## Server (cmd/server)

The server is an application that receives metrics from the agent, displays them in a browser, and stores them in the chosen storage (supports memory, file, postgresql).

### Compiling the Server

To compile the server, execute the following command:

`make build-server` 

### Using the Server

#### Server Parameters

-   `--address (or env var ADDRESS)`: The address on which the server will run (default :8080).
-   `--logs`: Logging level (info, debug).
-   `--store_interval (or env var STORE_INTERVAL)`: Used if storage = disk, specifying how frequently to write received data from memory to disk (default 300 seconds).
-   `--file_storage_path (or env var FILE_STORAGE_PATH)`: Used if storage = disk, specifying the path to the file for storing agent metrics (default ./tmp/metrics-db.json).
-   `--restore (or env var RESTORE)`: Whether to load data from storage during server initialization.
-   `--database_dsn (or env var DATABASE_DSN)`: Connection string for connecting to PostgreSQL.
-   `--key (or env var KEY)`: The key for verifying the signature of messages received from the agent.
  
## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](https://github.com/korovindenis/go-pc-info/blob/master/LICENSE.txt) file for details.

## Contributions and Feedback

Contributions to this project are welcome! If you have suggestions, find issues, or want to contribute new features or improvements, please feel free to open an issue or a pull request on GitHub.

Your feedback is valuable to us. If you have any questions or encounter any issues while using go-pc-metrics, please don't hesitate to reach out and open an issue on the GitHub repository. We appreciate your support in making this project better.
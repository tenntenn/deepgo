# deepgo [![PkgGoDev](https://pkg.go.dev/badge/github.com/tenntenn/deepgo)](https://pkg.go.dev/github.com/tenntenn/deepgo)

> **Warning**  
> `deepgo` is currently **experimental**. Breaking changes may occur.

`deepgo` provides various tools for deeper exploration of the Go ecosystem.

---

## MCP Server

### Tools

- **`proposal/review-meeting-minutes`**  
  Retrieves the latest Go proposal weekly meeting minutes from GitHub issue #33502.  
  **Parameters**  
  - **`limit`** (int): The maximum number of meeting minutes to fetch.  
  - **`since`** (string): Filter meeting minutes from a specific date/time (`YYYY-MM-DD hh:mm:ss`).

### Installation

Install `deepgomcp` with the following command:

```bash
go install github.com/tenntenn/deepgo/cmd/deepgomcp@latest
```

### Configuration

Below are example configurations for each platform. Adjust them as needed for your environment.

#### Mac

```json
{
  "mcpServers": {
    "deepgo": {
      "command": "deepgomcp",
      "args": [],
      "env": {}
    }
  }
}
```

#### Linux

```json
{
  "mcpServers": {
    "deepgo": {
      "command": "deepgomcp",
      "args": [],
      "env": {}
    }
  }
}
```

#### Windows (WSL)

```json
{
  "mcpServers": {
    "deepgo": {
      "command": "wsl",
      "args": [
        "bash",
        "-ic",
        "deepgomcp"
      ],
      "env": {}
    }
  }
}
```

---

## License

This project is licensed under the [MIT License](./LICENSE).

Contributions are always welcome! Please open issues or PRs for any bugs or enhancements.

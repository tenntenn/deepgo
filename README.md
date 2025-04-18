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
  - `limit` (int): The maximum number of meeting minutes to fetch
  - `since` (string): Filter meeting minutes from a specific date/time (`YYYY-MM-DD hh:mm:ss`)

- **`document/release-note`**
  Fetches the release note for the specified Go version from [go.dev/doc](https://go.dev/doc).  
  **Parameters**  
  - `version`* (string): The Go language version (e.g., `go1.24`)

- **`document/latest-go-version`**
  Retrieves the latest Go version from [https://go.dev/VERSION?m=text](https://go.dev/VERSION?m=text).  
  **Parameters**  
  - *(none)*

- **`style/modernize`**  
  Retrieves the source code of the modernize analyzer in gopls internal. All Go users must follow its rules.    
  **Parameters**  
  - *(none)*

- **`style/skeleton`**  
  The tool generates skeleton code of linter using golang.org/x/tools/go/analysis.Analyzer with github.com/gostaticanalysis/skeleton
  **Parameters**  
  - `kind` (string): The parameter kind (inspect,ssa,codegen,packages) is kind of skeleton code.)
  - `module` (string/required): The parameter module is a module path of skeleton code
  - `dst` (string/required): The dst parameter represents destination of files which must be absoluted path. When the dst parameter is empty string, the tool return generated code as txtar format string.

- **`util/copy-txtar`**  
  The tool copy files to given directory from txtar format string.
  **Parameters**  
  - `dir` (string/required): the dir parameter represents destination of files which must be absoluted path
  - `txtar` (string/required): the txtar parameter represents txtar fomrat string

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

# Fork-specific changes

This fork keeps Gin's standard-library JSON codec and the optional `go_json`
codec, while intentionally removing:

- BSON binding, rendering, negotiation, and the MongoDB Go Driver dependency.
- The `sonic` JSON build-tag implementation.
- The `jsoniter` JSON build-tag implementation.

These removals keep `go.mongodb.org/mongo-driver/v2`,
`github.com/bytedance/sonic`, `github.com/json-iterator/go`,
`github.com/modern-go/concurrent`, and
`github.com/twitchyliquid64/golang-asm` out of the module graph.

When merging upstream Gin changes, resolve conflicts by preserving these
removals. Validate the fork with:

```sh
go mod tidy
go test ./...
go test -tags nomsgpack ./...
go test -tags go_json ./...
```

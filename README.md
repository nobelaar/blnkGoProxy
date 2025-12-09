# BLNK Go Proxy

Lightweight proxy that adapts JavaScript/TypeScript clients to BLNK Core Banking to avoid large-integer issues. It lets clients send `precise_amount` as a string and converts it to an integer before forwarding the request to the core, preventing errors like `Do not know how to serialize a BigInt` or values being turned into exponential notation.

## What It Does
- Listens on `ProxyPort` and forwards everything to the host/port set in `config/config.go`.
- If the body is JSON and contains `precise_amount` as a string, it converts it to an integer and logs the conversion.
- If `precise_amount` is missing, already numeric, or cannot be converted, the request is forwarded unchanged (invalid cases are logged as warnings).

## Requirements
- Go >= 1.24 (see `go.mod`).
- Network connectivity to the configured BLNK instance.

## Configuration
Edit `config/config.go` to point to the core:
- `TargetHost`: BLNK host or IP.
- `TargetPort`: BLNK port.
- `ProxyPort`: port where the proxy listens (defaults to `5000`).

## Run
```bash
go run .
# or
go build -o blnkGoProxy && ./blnkGoProxy
```

Example from a JS/TS client (sending the amount as a string):
```bash
curl -X POST http://localhost:5000/transactions \
  -H "Content-Type: application/json" \
  -d '{"precise_amount":"123456789012345678", "account_id":"abc123"}'
```
The proxy converts `precise_amount` to an integer and forwards the request to BLNK at `TargetHost:TargetPort`. The core’s response is returned unchanged to the client.

## Notes and Limitations
- Only processes `application/json` bodies and the top-level `precise_amount` field.
- Conversion uses Go `int`; ensure values fit the target platform limits.
- Conversion errors do not stop the request: the original body is forwarded and a warning is printed to stdout.

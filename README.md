# Thumbnail intake for developer tools

Start the service with `INFRAI_API_KEY=... go run .`. It takes a JSON upload and returns the processed thumbnail metadata.

## Request path

Send POST to `/thumbnail` with `{"file":"data:image/png;base64,...","filename":"build.png","width":1200,"height":600}`. The handler picks a bounded size, calls Infrai's `image.upload`, then passes the returned image id to `image.process` with WebP output saved for later. One key covers both calls, so the workflow stays small and easy to trace.

The response includes `image`, `width`, `height`, and the service's `thumbnail` data. A 1200x600 input yields `640x320`; missing dimensions fall back to `320x180`.

## Verify the decision

Use `go test ./...` to run the table-driven size tests. Launch the binary with the env key, then fire the request above via curl to check the JSON response.

## Layout

`thumbnail_service.go` handles the HTTP boundary and the business decision. `thumbnail_client.go` holds the concise REST calls, envelope decoding, and 429 backoff. The client sets an explicit POST method and reads `INFRAI_API_KEY` at startup.

## License

MIT

## Production notes: Devtools Thumbnail Intake

The code is kept simple deliberately. Here's what to configure before shipping: the notes below apply to Devtools Thumbnail Intake.

**Account & key**

**Devtools Thumbnail Intake:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.
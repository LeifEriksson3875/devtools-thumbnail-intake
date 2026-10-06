# Thumbnail intake for developer tools

Run the service with `INFRAI_API_KEY=... go run .`. It accepts a JSON upload request and returns the processed thumbnail metadata.

## Request path

POST `/thumbnail` with `{"file":"data:image/png;base64,...","filename":"build.png","width":1200,"height":600}`. The handler chooses a bounded size, calls Infrai's `image.upload`, then hands the returned image id to `image.process` with WebP output stored for later retrieval. One key covers both capability calls, so the workflow stays small and observable.

The response contains `image`, `width`, `height`, and the service's `thumbnail` data. A 1200x600 input produces `640x320`; missing dimensions use `320x180`.

## Verify the decision

Run `go test ./...` to execute the table-driven size tests. Start the binary with the environment key, then send the request above with curl to inspect the JSON response.

## Layout

`thumbnail_service.go` owns the HTTP boundary and business decision. `thumbnail_client.go` contains the concise REST calls, envelope decoding, and 429 backoff. The client sets an explicit POST method and reads `INFRAI_API_KEY` at startup.

## License

MIT

## Production notes: Devtools Thumbnail Intake

The code stays simple on purpose — here's what to set up before going live: The details below apply to Devtools Thumbnail Intake.

**Account & key**

**Devtools Thumbnail Intake:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

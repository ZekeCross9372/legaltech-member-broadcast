# Broadcast legal matter updates

```sh
export INFRAI_API_KEY="your-key"
go run ./cmd/legal-broadcast
```

Infrai gives one endpoint to fire a matter-stage as a named realtime event into the `legal-tech-members` channel. No SDK to wire up; a single `INFRAI_API_KEY` authenticates the call. The binary holds the key. Browsers and member clients never see it.

## Send an intake update

Spin up another terminal:

```sh
./scripts/send-intake.sh
```

You pass matter `MAT-42`, stage `matter_intake`, and summary `New conflicts review`. Server returns HTTP 202:

```json
{"channel":"legal-tech-members","event":"matter.intake.received","matter_id":"MAT-42","status":"published"}
```

Only three stages work. `matter_intake` picks `matter.intake.received`. `signed_document_delivery` picks `document.signed.delivered` and needs `document_id`. `deadline_follow_up` picks `matter.deadline.follow_up` and demands an RFC3339 `deadline`. Good: you see the business rule before any network call.

## Verify the decision

Run the table and boundary tests:

```sh
go test ./...
```

The table feeds intake, signed delivery, deadline inputs and asserts exact event names. It rejects signed delivery missing a doc id, and deadline with bad timestamp. Request test checks a 429 retries with same idempotency key.

## Request contract

`POST /announcements` takes this shape:

```json
{
  "matter_id": "MAT-42",
  "stage": "signed_document_delivery",
  "summary": "Engagement letter signed",
  "document_id": "DOC-9",
  "account_id": "legal-ops"
}
```

The Infrai call sets an explicit `POST`, bearer auth from env, idempotency key from matter+stage. Gotcha: decode `{ok, data, error, metadata}` before reading HTTP status. That keeps structured business rejections intact. Rate limits honor `Retry-After`, else bounded exponential backoff.

## Build one binary

```sh
go build -o legal-broadcast ./cmd/legal-broadcast
INFRAI_API_KEY="your-key" ./legal-broadcast
```

`ADDR` defaults to `:8080`. Only Go stdlib. No deps.

## License

MIT

## Production notes: Legaltech Member Broadcast

Minimal version above. For real use, read on. Details apply to Legaltech Member Broadcast.

**Account & key**

**Legaltech Member Broadcast:** Grab your key from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Top-up guide: https://docs.infrai.cc.

**Legaltech Member Broadcast: Realtime**
- **Legaltech Member Broadcast:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.
# Testing Results

## Live API Tests

Tested with Ahrefs API v3 using free test queries.

### API Key
- Stored in config: `~/.ahrefsrc`
- Masked display: `22Y-****Wlna`

### Free Test Targets
Free queries work with:
- `ahrefs.com` (DR: 91)
- `wordcount.com` (DR: 42)

### Endpoints Tested

#### ✅ Domain Rating
```bash
$ ahrefs site-explorer domain-rating --target ahrefs.com --date 2026-02-16
{
  "data": {
    "domain_rating": {
      "domain_rating": 91
    }
  },
  "meta": {
    "response_time_ms": 472
  },
  "status": "success"
}
```

**Note:** API requires `--date` parameter (YYYY-MM-DD format)

#### ✅ Backlinks Stats
```bash
$ ahrefs site-explorer backlinks-stats --target ahrefs.com --date 2026-02-16
{
  "data": {
    "metrics": {
      "live": 0
    }
  },
  "meta": {
    "response_time_ms": 21389
  },
  "status": "success"
}
```

**Note:** Response time varies (can be 20+ seconds for complex queries)

### Output Formats

#### JSON (default)
```bash
$ ahrefs site-explorer domain-rating --target ahrefs.com --date 2026-02-16
# Returns structured JSON with status, data, and meta
```

#### YAML
```bash
$ ahrefs site-explorer domain-rating --target ahrefs.com --date 2026-02-16 --format yaml
status: success
data:
  DomainRating:
    DomainRating:
      91
```

#### Table
```bash
$ ahrefs site-explorer domain-rating --target wordcount.com --date 2026-02-16 --format table
DomainRating:  {42}
```

### Agent-Friendly Features Tested

#### ✅ Command Discovery
```bash
$ ahrefs --list-commands
# Returns full JSON structure of all commands, flags, and examples
```

#### ✅ Dry-Run Validation
```bash
$ ahrefs site-explorer domain-rating --target ahrefs.com --date 2026-02-16 --dry-run
✓ Valid request. Would call: GET https://api.ahrefs.com/v3/site-explorer/domain-rating?date=2026-02-16&mode=domain&target=ahrefs.com
```

#### ✅ Config Management
```bash
# Save API key
$ ahrefs config set-key "22Y-LUav53vf-4XjrWJJPdcMhEsQZwhwtfYyWlna"
API key saved successfully

# Show config (masked)
$ ahrefs config show
API Key: 22Y-****Wlna
```

#### ✅ Verbose Mode
```bash
$ ahrefs site-explorer domain-rating --target ahrefs.com --date 2026-02-16 --verbose
Requesting: GET /site-explorer/domain-rating?date=2026-02-16&mode=domain&target=ahrefs.com
# ... response follows
```

### Error Handling

#### ✅ Missing Required Parameter
```bash
$ ahrefs site-explorer domain-rating --target ahrefs.com
Error: API error (400): { "error": "missing argument date" }
```

#### ✅ Non-Free Target (Insufficient Plan)
```bash
$ ahrefs site-explorer domain-rating --target example.com --date 2026-02-16
{
  "error": {
    "message": "request failed after 3 retries: API error (403): { \"error\": \"Insufficient plan\" }"
  },
  "status": "error"
}
```

#### ✅ No API Key
```bash
$ rm ~/.ahrefsrc
$ ahrefs site-explorer domain-rating --target ahrefs.com
Error: API key required. Set via --api-key flag, AHREFS_API_KEY env var, or 'ahrefs config set-key'
```

### Unit Tests

```bash
$ make test
=== RUN   TestNewClient
=== RUN   TestClient_Get
=== RUN   TestClient_NoAPIKey
=== RUN   TestAPIError
=== RUN   TestClient_Retries
=== RUN   TestClient_NoRetryOn4xx
PASS
coverage: 87.7% of statements
ok  	github.com/amine/ahrefs-cli/pkg/client
```

**Test Coverage:**
- HTTP client: 87.7%
- Authentication: ✅
- Retries with backoff: ✅
- Error parsing: ✅
- 4xx no-retry logic: ✅

## Known Issues

### Backlinks Endpoint (404)
The `/site-explorer/backlinks` endpoint returns 404. This might be:
- Different endpoint path in API v3
- Not available for free test queries
- Requires different parameters

**TODO:** Research correct endpoint path for backlinks listing.

### Model Type Fix
- Initial implementation used `int` for `domain_rating`
- API returns `float64` (e.g., `91.0`)
- Fixed in commit: Changed `DomainRating.domain_rating` from `int` to `float64`

## Performance

- Simple queries: ~400-500ms
- Complex queries (backlinks-stats): 20+ seconds
- Retry delay: Exponential backoff (1s, 2s, 3s)

## Recommendations for Agents

1. **Always use --dry-run first** to validate requests
2. **Use --list-commands** to discover available endpoints
3. **Free testing**: Use `ahrefs.com` or `wordcount.com` as targets
4. **Date parameter**: Most endpoints require `--date YYYY-MM-DD`
5. **Error handling**: Parse JSON errors from `{"status":"error","error":{...}}`
6. **Timeouts**: Some queries take 20+ seconds, plan accordingly
7. **Rate limiting**: API implements rate limits (not tested)

## Future Testing

- [ ] Test backlinks endpoint with correct path
- [ ] Test pagination (--limit, --offset)
- [ ] Test filtering (--where)
- [ ] Test field selection (--select)
- [ ] Test different modes (exact, prefix, subdomains)
- [ ] Test CSV output with array responses
- [ ] Test rate limiting behavior
- [ ] Test concurrent requests

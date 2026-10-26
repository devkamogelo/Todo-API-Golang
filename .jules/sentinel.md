## 2024-10-26 - Information Disclosure in HTTP 500 Responses
**Vulnerability:** The application was directly exposing internal error details (like database errors or internal context errors) in HTTP 500 response bodies by sending `err.Error()` in JSON.
**Learning:** Returning `err.Error()` to clients in 500 status responses can leak sensitive system information (such as database structure, failed paths, or internal logic). This fails the "Fail securely" principle.
**Prevention:** Always log the actual error internally (e.g., using `c.Error(err)`) and return a generic error message (like "Internal server error") to the client for 500 status codes. Also, always ensure there is a `return` statement after `c.JSON()` in error paths to prevent fall-through.

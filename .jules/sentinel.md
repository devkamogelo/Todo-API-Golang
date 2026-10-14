## 2024-05-18 - [Error Message Information Disclosure]
**Vulnerability:** Internal server errors in HTTP handlers returned `err.Error()` directly to the client in JSON responses (e.g., `Failed to hash password " + err.Error()`).
**Learning:** Detailed error messages can leak sensitive internal information such as stack traces, database schema details, or underlying library versions to potential attackers.
**Prevention:** Implement secure error handling by returning generic error messages (e.g., "An internal server error occurred") for unexpected faults. Log the detailed error internally for debugging instead of exposing it via the API.

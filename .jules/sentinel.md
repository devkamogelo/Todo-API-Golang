## 2024-05-18 - Prevent Error Detail Leakage in API Responses
**Vulnerability:** Information disclosure via detailed error messages (e.g. database constraint errors, bcrypt errors, JWT signing errors) directly in HTTP 500 Internal Server Error JSON responses.
**Learning:** Detailed error messages can leak internal database structure (e.g., column names, constraint names) or internal application logic which can be used by an attacker to map the backend and craft further attacks.
**Prevention:** Catch and log errors securely on the backend (e.g., using `c.Error(err)`) and only return a generic "Internal server error" message to the user. Always follow with a `return` statement in Gin handlers to prevent fail-open execution flows.

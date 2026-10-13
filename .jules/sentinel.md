## 2024-10-24 - Missing Gin Returns & Information Leakage
**Vulnerability:** Found multiple handlers that called `c.JSON(http.StatusInternalServerError, ...)` without returning immediately after, leading to potential panics or DoS due to multiple response writes. Also found the application leaking raw `err.Error()` to clients on internal server errors, exposing database and internal logic.
**Learning:** Handlers using Gin must `return` after writing a terminal response (like `c.JSON`) if execution shouldn't continue. Unhandled errors exposed raw application internals.
**Prevention:** Ensure every `c.JSON` that signifies an early return (error cases) is followed by a `return` statement. Use generic error messages like "An internal server error occurred" instead of raw errors.

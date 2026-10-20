## 2024-03-24 - Fix LogoutHandler Nil-Pointer Dereference
**Vulnerability:** A critical nil-pointer dereference vulnerability (Denial of Service) was found in `LogoutHandler` where it attempted to access `claims.ID` when token parsing failed.
**Learning:** The check `if err != nil` was incorrectly used instead of `if err == nil`. When token parsing failed, `claims` was nil, causing a panic when `claims.ID` was accessed.
**Prevention:** Always ensure that pointers are not nil before dereferencing them. Double check error handling logic to ensure the success path is taken when `err == nil`.
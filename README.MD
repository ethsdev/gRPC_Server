# LocatingRecordsInLogs-gRPC

A distributed gRPC-based application for locating requested records in log files via a Lambda-backed HTTP API.

---
**Platform:** Windows 10  
**Go Version:** 1.19.2  
**IDE:** IntelliJ IDEA Ultimate

---

## 🚀 Project Overview

This project consists of a gRPC client-server architecture where:

- The **gRPC client** sends a request with a specific `time` and `deltaTime`.
- The **gRPC server** receives the request and queries an AWS API Gateway endpoint.
- The endpoint routes the request to a Lambda function that checks for log records within the specified time range.

### ✅ Response Behavior

- If logs **are present**, a 200 response is returned with an MD5 hash of the logs.
- If logs **are not found**, a 400-level response is returned with a "not found" message.

---

## 🔧 Refactored Architecture

> ⚠️ **Major structural changes implemented:**

- Extracted API call logic into a reusable function.
- Added structured error handling (replaced `os.Exit` and raw `fmt.Print`).
- Used `http.Client` with timeouts.
- Switched from deprecated `ioutil` to `io` for reading responses.
- Embedded `UnimplementedTestApiServer` for forward compatibility.
- Added constructor method `NewTestApiServer()` for better code organization and testability.

---
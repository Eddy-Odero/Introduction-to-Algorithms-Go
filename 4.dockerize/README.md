# ASCII Art Web - Dockerized

A production-ready, containerized web application written in Go that converts user input into ASCII art using different banner styles. The project leverages a highly optimized multi-stage Docker configuration to satisfy advanced garbage collection requirements.

---

## Features
- Generate ASCII art from user input via a clean web interface.
- Supports multiple styles: Standard, Shadow, and Thinkertoy.
- Statically compiled Go binary for maximum speed and minimal memory footprint.
- Staged deployment layer resulting in an extremely small container footprint (~25MB).
- Imbedded object metadata labels for audit verification.

---

## Technologies Used
- **Go** (Standard Packages Only: `net/http`, `html/template`, `fmt`, `log`)
- **Docker** & **Docker CLI Client Utilities**
- **HTML / CSS**

---

## How to Run (Docker Deployment Workflow)

Follow these steps to build the image and spin up the container exactly as required by the auditor checklist.

### Step 1: Clone the Repository
```bash
git clone https://github.com
cd Introduction-to-Algorithms-Go/4.dockerize
```

### Step 2: Build the Container Image
Compile the isolated application layer and apply embedded structural configuration instructions:
```bash
docker image build -f Dockerfile -t ascii-art-web:1.0 .
```

### Step 3: Run the Container Service
Launch the containerized server as a background daemon mapping port 8080 from the container to your local machine:
```bash
docker container run -p 8080:8080 --detach --name ascii-art-container ascii-art-web:1.0
```

### Step 4: Access the Application
Open your web browser and navigate to:
```text
http://localhost:8080
```

---

## Auditor Verification & Test Requirements

Execute these commands to audit the container infrastructure and ensure compliance with all project rules.

### Test 1: Verify Active Running State
Check that the container is active and securely listening on port 8080:
```bash
docker ps -a
```
*Expected Output: Look for `ascii-art-container` showing a status of `Up X seconds` or `Up X minutes`.*

### Test 2: Verify Applied Metadata
Ensure custom object metadata strings are applied directly to the image layers:
```bash
docker inspect --format='{{json .Config.Labels}}' ascii-art-web:1.0
```
*Expected Output: A JSON payload containing the maintainer name (`Eddy-Odero`), version, and project description.*

### Test 3: Inspect Container File System
Verify that unused files (like source code files and Go compilers) were correctly left behind during the garbage collection process:
```bash
docker exec -it ascii-art-container ls -la
```
*Expected Output: A minimal runtime environment containing only the execution binary (`ascii-art-web-server`), the `template` layout directory, and the `banners` folder.*

### Test 4: Verify Image Size Optimization (Garbage Collection)
Verify that the production layer strips away the 800MB Go compiler tools:
```bash
docker images ascii-art-web:1.0
```
*Expected Output: The total image size column must reflect a highly-optimized profile between 15MB and 35MB.*

---

## Application Error Handling Matrix

You can manually trigger and verify our strict HTTP state restrictions directly through your browser or CLI utilities:

1. **404 Not Found Handling**
   - **Test:** Go to `http://localhost:8080/unknown-route`
   - **Expected Result:** Server responds with a clean `404 Page Not Found` message.
2. **405 Method Not Allowed Handling**
   - **Test:** Run `curl -X GET http://localhost:8080/ascii-art`
   - **Expected Result:** Server rejects the direct request with a `405 Method Not Allowed` payload.
3. **500 Internal Server Error Isolation**
   - **Test:** Intentionally break a path tracking variable inside the assets directory.
   - **Expected Result:** System gracefully isolates the defect and serves an `Internal Server Error` state without crashing the main container background process.

---

## Running Unit Tests Locally
If you have a local Go toolchain configured, you can execute the raw algorithmic unit test blocks using:
```bash
go test ./...
```

**Author:** Eddy-Odero  
**Project Track:** Dockerize/Infrastructure Optimization

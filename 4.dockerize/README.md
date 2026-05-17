# ASCII Art Web

A web application written in Go that converts user input into ASCII art using different banner styles.

---

## Features

* Generate ASCII art from user input
* Supports multiple styles:

  * Standard
  * Shadow
  * Thinkertoy
* Clean web interface
* Error handling for invalid input and server issues

---

## Technologies Used

* Go
* HTML/CSS 
* net/http package

---

##  How to Run

1. Clone the repository:

```bash
git clone https://github.com/Eddy-Odero/Introduction-to-Algorithms-Go/tree/main/3.ascii-art-web
cd ascii-art-web
```

2. Run the server:

```bash
go run .
```

3. Open your browser:

```
http://localhost:8080
```

---

##  Project Structure

```
ascii-art-web/
│
├── main.go
├── ascii/
│   ├── loader.go
│   ├── engine.go
│   └── engine_test.go
│
├── templates/
│   └── index.html
│
├── banners/
│   ├── standard.txt
│   ├── shadow.txt
│   └── thinkertoy.txt
```

---

##  Running Tests

Run all tests using:

```bash
go test ./...
```

---

##  Error Handling Tests

### 1. 404 Not Found

**Test:**

```
http://localhost:8080/unknown
```

**Expected Result:**

* Server returns 404 Not Found

---

### 2. 400 Bad Request (Empty Input)

**Test:**

1. Remove `required` attribute from the textarea in `index.html`
2. Submit an empty form

**Expected Result:**

* Error message displayed OR HTTP 400 response

---

### 3. 500 Internal Server Error (Missing Banner File)

**Test:**

1. Rename or delete a banner file:

```
banners/standard.txt
```

2. Submit a request

**Expected Result:**

* Server returns Internal Server Error (500)

---

### 4. Invalid Banner Input

**Test:**

* Modify form manually (via browser DevTools) and send an invalid banner value

**Expected Result:**

* Error message OR HTTP 400 response

---

### 5. Invalid HTTP Method

**Test:**

```
http://localhost:8080/ascii-art
```

**Expected Result:**

* Method Not Allowed error

---

## Normal Usage Test

1. Enter text in the input field
2. Select a banner style
3. Click **Generate ASCII**

**Expected Result:**

* ASCII art is displayed correctly on the same page

---

## Notes

* The server handles both user input validation and internal errors
* ASCII banners are loaded from text files
* Output is rendered using Go templates

---

##  Author

* Eddy-Odero

---

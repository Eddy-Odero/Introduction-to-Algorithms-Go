#go-reloaded (Go Project)

A text processing program written in Go that applies a set of transformations and formatting rules to input text from a file, writes the result to an output file, and prints the final result to the terminal.

---

##  Features

This program processes text with the following capabilities:

###  Text Modifiers
- `(up)` → converts previous word to UPPERCASE
- `(up, N)` → converts previous N words to UPPERCASE
- `(low)` → converts previous word to lowercase
- `(low, N)` → converts previous N words to lowercase
- `(cap)` → capitalizes previous word
- `(cap, N)` → capitalizes previous N words

### Number Conversions
- `(hex)` → converts previous hexadecimal number to decimal
- `(bin)` → converts previous binary number to decimal

###  Articles Fix
- Automatically corrects:
  - `a apple` → `an apple`
  - `A elephant` → `An elephant`

###  Punctuation Formatting
- Removes unwanted spaces before punctuation:
  - `hello , world` → `hello, world`
- Handles grouped punctuation:
  - `...` stays intact
  - `!?` stays intact

###  Quotes Handling
- Cleans spaces inside single quotes:
  - `' hello world '` → `'hello world'`

---

##  Project Structure

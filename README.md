
# a0 Lang

a0 is a simple and expressive programming language with fun aliases for keywords and a clean syntax. It supports functions, variables, control flow, logical expressions, and immediate interpretation.

---

## Heads up!

This was a project I made a good while ago and might have bugs I don't know about, but in the end
it was made for fun

---

## Features

- Function declarations with multiple aliases (`func`, `fun`, `fn`, `funky`, `def`)  
- Variable declarations with aliases (`var`, `val`, `let`, `define`)  
- Constants with `const`  
- Control flow with `if`, `for`, `while` and fun synonyms like `loop`, `forever`  
- Logical operators with aliases (`and` / `&&` / `plus`, `or` / `||` / `perhaps`)  
- Unicode keyword support (e.g., `❓` for `if`)  
- Prints a clear, human-readable AST for debugging  
- Simple interpreter to run your a0 programs  

---

## Getting Started

### Prerequisites

- Go 1.21 or later  

### Build

Clone the repository and build:

```bash
git clone https://github.com/Mstr0A/a0-lang.git
cd a0-lang
go build -o a0
````

### Usage

Run your source file like this:

```bash
./a0 path/to/yourfile.a0
```

Flags:

* `-tokens` — Print token list and exit
* `-ast` — Print AST and exit

Example:

```bash
./a0 -ast example.a0
```

---

## Sample Code

```a0
funky greet(name, num) {
    print("Hello, ", name, num, "!")
}

funky factorial(num) {
    val result = 1
    while (num > 1) {
        result = result * num
        num = num - 1
    }
    return result
}

val counter = 3
while (counter > 0) {
    greet("User #", counter)

    print("Factorial of ", counter, " is ", factorial(counter))
    print("")

    counter = counter - 1
}

print("Done")
```

---

## Keywords and Aliases

| Keyword(s)                          | Meaning              |               |            |
| ----------------------------------- | -------------------- | ------------- | ---------- |
| `func`, `fun`, `fn`, `funky`, `def` | Function declaration |               |            |
| `var`, `val`, `let`, `define`       | Variable declaration |               |            |
| `const`                             | Constant declaration |               |            |
| `if`, `❓`                          | Conditional          |               |            |
| `for`                               | For loop             |               |            |
| `while`, `loop`, `forever`          | While loop           |               |            |
| `and`, `plus`                       | Logical AND          |               |            |
| `or`, `perhaps`                     | Logical OR           |               |            |
| `not`, `!`                          | Logical NOT          |               |            |

---

## License

MIT License — use it as you want.

---

## Contributing

Contributions, issues, and ideas are welcome!

---

```

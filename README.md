# Templang

_Name pending_

Templang is a dynamically typed scripting language made created for me to learn
about how programming languages are made.

## Language features

- [x] C-like syntax
- [x] mutable variables
- [x] arithmetic operations
- [x] boolean operations
- [x] `print`-statement (temporary until a print function is implemented in std)
- [x] Infix expressions
- [x] prefix expression

### upcomming freatures

- [ ] functions
- [ ] first class function
- [ ] higher order function
- [ ] closures
- [ ] if/else control flow
- [ ] iteration (for statement)
- [ ] complex data structures (array and map)
- [ ] lexical scoping
- [ ] _

## about

everything is written by hand, without the use of any other library than the go std lib.
This is mosty based on the books written by Robert Nystrom and Thorsten Ball (links at the bottom).

### Lexer

The lexer/scanner/tokenizer is the first part of of the interpreter. This is responsible for
splitting the text into meaningful peices if data. This data is called a "Token". A token represents
a substring of the source text, and has some extra fields added to help us further down the road. This includes
the part of source code that resulted in the token being parsed, a type that we can use for comparing/expecting
tokens while parsing and line/col data to make better errors.

### Parser

The parser takes list of tokens and creates an
[abstract syntax tree](https://en.wikipedia.org/wiki/Abstract_syntax_tree) (AST)
based on the type/value of tokens. For most of this paring, we use a
[Pratt parsing](https://en.wikipedia.org/wiki/Operator-precedence_parser#Pratt_parsing)
algorith.

### Evaluator/interpreter

At last, we evaluate the AST based on the semantics defined on the interpreter.

## sources

- [Crafting interpreters by Robert Nystrom](https://craftinginterpreters.com)
- [Pratt parsers: expression parsing made easy](https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/)
- [Writing An Interpreter In Go by Thorsten Ball](https://interpreterbook.com)
- [Writing A Compiler In Go by Thorsten Ball](https://compilerbook.com)
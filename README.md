# Templang

_Name pending_

Templang is a dynamically typed scripting language made for me to learn
about how programming languages are made.

## (Current) Language features

- [x] C-like syntax
- [x] Line comments
- [x] mutable variables
- [x] arithmetic operations
- [x] boolean operations
- [x] `print`-statement (temporary until a print function is implemented in std)
- [x] Infix expressions
- [x] prefix expression
- [x] lexical scoping
- [x] if/else control flow
- [x] functions
- [x] first class function
- [x] higher order function
- [x] closures
- [x] return statements
- [x] attach REPL to an executed program
- [x] iteration (for loops or some form of iterator implementation). Currently supports the follow values to iterate:
  - number - does n iterations where n=number
  - string - loops through each char in a string. should handle UTF-8 correctly
  - boolean - infinite loop on true, skip on false
  - list - loop through items in a list from start to end
  - none - default to boolean=true
- [x] "list" value for ordered lists of values
  - Indexing into list
  - Iteration over list
  - Push/pop for lists
  - Assign value at index
- [x] Add "map"-object for key-value pairs where all primitive values are valid keys
  - Indexing into map
  - Iteration over map (by key)
  - Assign value at key
- [x] Builtin functions
  - push - add element(s) to end of a list
  - pop - remove the last element of list
  - len - return length of list/map/string
  - str - return the value as its string representation
- [x] module system with importing from std lib/another file. requires:
  - language support for accessing members of namespaces (syntax, parsing and resolving)
  - expanding the internal typing to support multiple sources
  - import resoultion
- [x] some form of std lib implemented with the language
  - [ ] http
  - [ ] math
  - [ ] fmt
  - [ ] iteration lib - iter
  - [ ] ...

### upcoming features / TODOs

- [ ] \[IDEA\] add range/slice operator for indexing and loops
- [ ] allow for pull iteration with `iteratorObj` by exposing the internal `next()` and `done()` methods as properies of iteratorObj
  - syntax could be be something like `expr -> expr`
  - should support ranging positive and negative direction
  - should only be valid when used in each stmts.
  - could also just be a builtin function that creates an iterator for the range
- [ ] error messages on runtime errors that help you identify your error by pointing to the error in source code
- [ ] async primitive of some sort.
- [ ] formatted strings with print statment

## about

everything is written by hand, without the use of any other library than the go std lib.
This is mosty based on [Crafting interpreters] by Robert Nystrom, [Writing An Interpreter In Go] and [Writing A Compiler In Go] Thorsten Ball.

### Lexer

The lexer/scanner/tokenizer is the first part of of the interpreter. This is responsible for
splitting the text into meaningful peices if data. This data is called a "Token". A token represents
a substring of the source text, and has some extra fields added to help us further down the road. This includes
the part of source code that resulted in the token being parsed, a type that we can use for comparing/expecting
tokens while parsing and line/col data to make better errors.

### Parser

The parser takes list of tokens and creates an
[abstract syntax tree] (AST)
based on the type/value of tokens. For most of this paring, we use a
[Pratt parsing] algorith.

### Evaluator/interpreter

At last, we evaluate the AST based on the semantics defined on the interpreter. The interpreter is
whats called a "tree-walking-interpreter" where it works directly on the AST, instead of making
istead of a more efficient datastructure. The reason for this is simplicity.

## Notes

There are some blurred lines on where the responsibility lies when it comes to the interpreter.
The main idea is that the lexer and parser are responsible for making sure that the program
is grammatically correct, but they should not care about the semantics. an example of this with english:

- timmy drives a car.

This is a valid sentence in english following the structure of `subject-verb-object`.
The parser defines this structure, and will accept anything, as long as the structure is correct.

- emil dances a motorcycle.

This sentence follows the same structure as the previous one, but it doesn't make sense. The parser
doesn't care about that, as long as the words are the correct type, in the correct order, it is happy.

The meaning of the words is decided by the interpreter. The interpreter will evaluate the sentences
and see if the words actually mean anything.

<!-- ## sources -->

[Pratt parsing]: https://en.wikipedia.org/wiki/Operator-precedence_parser#Pratt_parsing
[abstract syntax tree]: https://en.wikipedia.org/wiki/Abstract_syntax_tree
[Crafting interpreters]: https://craftinginterpreters.com
[Pratt parsers: expression parsing made easy]: https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/
[Writing An Interpreter In Go]: https://interpreterbook.com
[Writing A Compiler In Go]: https://compilerbook.com

### CISP - CSS in Server Pages
A simple css like language to write simple server side scripts.
like php, but much much much worse.

## Syntax

id = ([a-zA-Z_][a-zA-Z0-9_\-]*) | (-- + id)
string = ('"' + .* + '"') | ("'" + .* + "'")
number = ([0-9]* + (. + [0-9]+)?)
boolean = true | false
value = string | number | boolean | function_call
unary_operator = ! | ~
operator = + | - | * | / | % | ^ | = | == | != | > | >= | < | <= | && | ||
expression = value | unary_operator + expression
    | function_call
    | (expression + operator + expression)
    | '(' + expression + ')'
function_call = id + '(' + function_parameters + ')'
function_parameters = expression + (, + expression)* + (,)*

selector = (. | #)? + id + ([id=value])* + (, + selector)*

at_rule = @ + id + expression + (delcaration_block | ;)
rule = (selector + declaration_block) | at_rule
delcaration_block = { + (declaration | rule)* + }
declaration = id + : + value;

program = (rule | at_rule)*

## Progress

- [x] Lexer
- [ ] Parser
    - [x] Selector
    - [x] Declaration
    - [x] Declaration Block
    - [x] Rule
    - [x] At Rule
    - [x] Program
    - [ ] Value
        - [x] String
        - [x] Number
        - [x] Boolean
        - [x] Function Call
        - [ ] Unary Operator
        - [ ] Expression
- [ ] Interpreter
- [ ] Compiler

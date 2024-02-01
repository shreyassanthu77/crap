### CISP - CSS Is a Programming language
CISP is a general purpose programming language inspired by CSS.
It is a functional language with a focus on simplicity and readability (if I say so myself).

## Getting Started
- No Comments, who needs them anyway. just write self-documenting code.
### Functions everywhere
- CSS rules you are used to are now functions.
- The Property Declarations are now function calls with the property name as the function name.
- The property value is now the function parameter. (multiple parameters are separated by spaces)
- If you want to call a function in the value of a property you can use `function_name()` syntax.
```css
thisIsAFunction {
    print: "Hello World!";
    add: 1 2;
    print: add(1, 2);
}
```
### Variables
- Any pproperty starting with `--` is a variable.
- Variables are scoped to the selector they are defined in and all child selectors (yes nested selectors are a thing).
- Variables can be used in any property value by using the `var` function
or by prefixing the variable name with a `$`
```css
main {
    --a: 1;
    --b: 2;
    --c: add(var(--a), $b);
    print: $c;
}
```
### Nesting
- Selectors can be nested inside each other (just like in modern CSS).
- Scopes work like you'd expect them in other languages.
```css
main {
    nested {
        print: "Hello World!";
    }

    nested: ();
}
```

### Function Parameters
- Function parameters can be declared as attributes on the selector.
- The value of the attribute is the default value of the parameter.
- If the parameter is not passed to the function the default value is used.
- If you don't want to pass a parameter you can use `()` as the value.
```css
someFunction[parameter=1] {
    print: $parameter;
}

main {
    someFunction: 2;
    someFunction: ();
}
```
prints:
```
2
1
```

## Examples

### Hello World
```css
main {
    print: "Hello World!";
}
```

### Factorial
```css
factorial[n] {
    @if $n == 0 || $n == 1 {
        @return 1;
    }

    @return $n * factorial($n - 1);
}
```

## Progress

- [x] Lexer
- [x] Parser
- [ ] Interpreter
- [ ] Compiler

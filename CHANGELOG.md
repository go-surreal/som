# Change Log

## 0.1.0 - First public release

- Fixed delete operations failing
- Delete, relate, and let queries now return `Promise<undefined>`
    - This makes the typings for arrays returned by chaning more accurate

## 0.0.5

- Added support for stateless queries
    - Use the new `CirqlStateless` class
    - Same API as the main stateful `Cirql` class
- Added query function for LET statement
- Added query function for IF ELSE statement
- Refactored authentication
    - Now supports scope, namespace, database, and token authentication
- Refactored AND & OR behavior
    - Now allows for more combinations than before

## 0.0.4
- Allow array add and remove in update queries

## 0.0.3
- Added support for query writers
- Added retry functionality

## 0.0.2
- Refactored some functions

## 0.0.1
- Basic functionality only

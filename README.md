## Running
```sh
go run . transactions.csv transactions.out # writes to file
```

or

```sh
go run . transactions.csv # writes to standard output
```

## TODO
* Better handling of invalid input/formulas (right now it mostly works for happy path)
* Implement RollbackWrapper parser and remove rollbacks from Map* parsers (workaround for parser.Sequence* bugs)
* Figure out better way to do infix ops (nested switches are really not pretty)

## Ideas
* Implement CalculatedValue for InvalidValue type, perhaps with trace to where the evaluation first failed
* Add inProgress flag to evalCell, to catch cycles in evaluation


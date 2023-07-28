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

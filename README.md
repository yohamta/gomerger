# gomerger

gomerger is a tool that merges all Go files in a project into a single Go file except for `*_test.go` files.

## Usage

```
go install github.com/yohamta/gomerger
go run merger.go /path/to/your/project
```

Then the output will be written to `merged_project.go` in the current directory.

## Notes

-   Test files (ending with `_test.go`) are excluded from the merge.
-   Currently, it may not create a valid Go file in all cases.

## Contributing

Any type of contributions are very welcome!

## License

MIT License

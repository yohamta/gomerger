# gomerger

gomerger is a tool that merges all Go files in a project into a single Go file.

## Features

-   Recursively scans a project directory for Go files
-   Merges all non-test Go files into a single output file
-   Preserves package structure and imports
-   Adds a main function if one doesn't exist in the original files

## Usage

```
go install github.com/yohamta/gomerger
go run merger.go /path/to/your/project
```

Then the output will be written to `merged_project.go` in the current directory.

## Notes

-   Test files (ending with `_test.go`) are excluded from the merge.
-   The tool preserves the original package names as comments in the merged file.
-   If no main function is found in the original files, an empty main function is added to the merged file.

## Limitations

-   Currently, it may not create a valid Go file in all cases.

## Contributing

Any type of contributions are very welcome!

## License

MIT License

Copyright (c) 2024 Yota Hamada

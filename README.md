# flatar

`flatar` is flattens nested directories and creates tar archives for each subdirectory

## Features

- Flatten nested directory structures
- Create tar or zip archives for flattened directories
- Flexible directory processing
- Handles naming conflicts when moving files

## Install

```bash
go install github.com/yashikota/flatar@latest
```

## Usage

### Basic Usage

Process subdirectories in the current directory

```bash
flatar
```

Process a specific directory

```bash
flatar /path/to/directory
```

### Command-line Options

```bash
flatar [-a] [-z] [-d] [<root_directory>]
```

- `-a`: Create tar archive after flattening
- `-z`: Create zip archive after flattening
- `-d`: Delete original directory after processing

### Error Handling

The tool provides informative messages in the following cases

- When the specified directory does not exist
- When the specified path is not a directory
- When no subdirectories are found to process

### Example

Given a directory structure

```txt
dir/
├── project1/
│   ├── fileA.txt
│   └── subdir/
│       └── fileA.txt
└── project2/
    └── fileB.txt
```

After running `flatar`

```txt
dir.tar/
├── fileA.txt
├── fileA_subdir.txt
├── fileB.txt
```

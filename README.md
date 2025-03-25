# flatar

`flatar` is flattens nested directories and creates tar archives for each subdirectory

## Features

- Flatten nested directory structures
- Create tar archives for flattened directories
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

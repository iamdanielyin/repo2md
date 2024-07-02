# repo2md

`repo2md` is a command-line tool for generating Markdown files that represent the directory structure and code content of a specified Git repository or local directory.

## Installation

Make sure you have Go installed on your system. Then, install `repo2md` using the following command:

```bash
go install github.com/iamdanielyin/repo2md@latest
```

## Usage

### Basic Usage

Generate a Markdown file for a local directory:

```bash
repo2md /path/to/your/local/repo
```

Generate a Markdown file for a remote Git repository:

```bash
repo2md https://github.com/username/repo.git
```

By default, the Markdown file will be named `repo_structure.md`.

### Options

- `-o <output-file>`: Specify the output Markdown file name (default is `repo_structure.md`).
- `-h`, `--help`: Display help information.

## Examples

Assume you have a local Git repository `/path/to/your/local/repo`. Use the following command to generate a Markdown file:

```bash
repo2md /path/to/your/local/repo
```

The generated Markdown file will include the directory structure and content of all code files in the repository.

## Notes

- The generated Markdown file includes the directory structure and content of all code files in the Git repository or local directory. Make sure you have read permissions for these files.
- If the Git repository contains a large number of files or large files, the generated Markdown file may be large.

## Support

If you encounter any issues or have any suggestions while using `repo2md`, please file an issue on [GitHub Issues](https://github.com/iamdanielyin/repo2md/issues).

---

This README provides comprehensive installation instructions, usage details, options, examples, notes, and support information to help users understand how to install and use the `repo2md` tool.

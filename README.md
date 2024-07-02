# repo2md

`repo2md` 是一个用于生成本地Git仓库目录结构和代码内容的命令行工具。它将指定的Git仓库或本地目录的文件结构输出为Markdown格式，并可以包含指定文件的代码内容。

## 安装

首先，确保你的系统已经安装了Go语言环境。然后使用以下命令安装 `repo2md`：

```bash
go install github.com/iamdanielyin/repo2md@latest
```

## 使用方法

### 基本用法

生成本地目录的Markdown文件：

```bash
repo2md /path/to/your/local/repo
```

生成远程Git仓库的Markdown文件：

```bash
repo2md https://github.com/username/repo.git
```

Markdown文件默认输出为 `repo_structure.md`。

### 选项

- `-o <output-file>`: 指定输出的Markdown文件名，默认为 `repo_structure.md`。
- `-h`, `--help`: 显示帮助信息。

## 示例

假设我们有一个本地的Git仓库 `/path/to/your/local/repo`，使用以下命令生成Markdown文件：

```bash
repo2md /path/to/your/local/repo
```

生成的Markdown文件将包含该仓库的目录结构和所有代码文件的内容。

## 注意事项

- 生成的Markdown文件将包含Git仓库的目录结构和所有代码文件的内容，确保你有权限读取这些文件。
- 如果Git仓库包含大量文件或者文件较大，生成的Markdown文件可能会比较大。

## 支持

如果你在使用过程中遇到问题或者有任何建议，请在 [GitHub Issues](https://github.com/iamdanielyin/repo2md/issues) 中提出。

---

这个README文件提供了完整的安装说明、使用方法、选项说明、示例、注意事项和支持信息，帮助用户了解如何安装和使用 `repo2md` 工具。
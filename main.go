package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/sabhiram/go-gitignore"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: repo2md <repo-path or repo-url>")
		return
	}

	input := os.Args[1]
	mdFile := "repo_structure.md"
	var repoPath string
	var err error

	// Check if the input is a local directory or a remote URL
	if isLocalDir(input) {
		repoPath = input
	} else {
		repoPath, err = cloneRepo(input)
		if err != nil {
			fmt.Println("Error cloning repository:", err)
			return
		}
		defer os.RemoveAll(repoPath) // Clean up the cloned repo after processing
	}

	// Load .gitignore if it exists
	gitignorePath := filepath.Join(repoPath, ".gitignore")
	var gitignorePatterns *ignore.GitIgnore
	if _, err := os.Stat(gitignorePath); err == nil {
		gitignorePatterns = loadGitignore(gitignorePath)
	} else {
		// If no .gitignore file, create an empty gitignore pattern
		gitignorePatterns = ignore.CompileIgnoreLines()
	}

	structure, content, err := generateMarkdown(repoPath, gitignorePatterns)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	mdContent := "# Code Structure\n\n" + structure + "\n# Code Content\n\n" + content
	err = os.WriteFile(mdFile, []byte(mdContent), 0644)
	if err != nil {
		fmt.Println("Error writing markdown file:", err)
		return
	}

	fmt.Println("Markdown file generated:", mdFile)
}

func isLocalDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func cloneRepo(url string) (string, error) {
	tempDir, err := os.MkdirTemp("", "repo")
	if err != nil {
		return "", err
	}

	_, err = git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		return "", err
	}

	return tempDir, nil
}

func loadGitignore(path string) *ignore.GitIgnore {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening .gitignore file:", err)
		return nil
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading .gitignore file:", err)
	}

	// Always ignore .git and LICENSE files
	patterns = append(patterns, ".git", ".gitignore", "LICENSE")

	return ignore.CompileIgnoreLines(patterns...)
}

func generateMarkdown(repoPath string, gitignorePatterns *ignore.GitIgnore) (string, string, error) {
	var structureBuilder, contentBuilder strings.Builder

	// Store directories and files separately
	dirMap := make(map[string][]string)

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return err
		}

		// Skip ignored files and directories
		if gitignorePatterns != nil && gitignorePatterns.MatchesPath(relativePath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			if relativePath != "." {
				dirMap[filepath.Dir(relativePath)] = append(dirMap[filepath.Dir(relativePath)], filepath.Base(relativePath)+"/")
			}
		} else {
			dirMap[filepath.Dir(relativePath)] = append(dirMap[filepath.Dir(relativePath)], filepath.Base(relativePath))
		}

		return nil
	})

	if err != nil {
		return "", "", err
	}

	// Sort directories and files
	for dir := range dirMap {
		sort.SliceStable(dirMap[dir], func(i, j int) bool {
			// Directories come first, then files, both in lexicographical order
			if strings.HasSuffix(dirMap[dir][i], "/") && !strings.HasSuffix(dirMap[dir][j], "/") {
				return true
			}
			if !strings.HasSuffix(dirMap[dir][i], "/") && strings.HasSuffix(dirMap[dir][j], "/") {
				return false
			}
			return dirMap[dir][i] < dirMap[dir][j]
		})
	}

	// Generate markdown for structure
	var generateStructure func(string, int)
	generateStructure = func(dir string, indent int) {
		for _, entry := range dirMap[dir] {
			structureBuilder.WriteString(strings.Repeat("    ", indent) + "- " + entry + "\n")
			if strings.HasSuffix(entry, "/") {
				generateStructure(filepath.Join(dir, strings.TrimSuffix(entry, "/")), indent+1)
			}
		}
	}

	generateStructure(".", 0)

	// Generate markdown for content
	err = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel(repoPath, path)
		if err != nil {
			return err
		}

		// Skip ignored files and directories
		if gitignorePatterns != nil && gitignorePatterns.MatchesPath(relativePath) {
			return nil
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if len(filepath.Ext(relativePath)) > 1 {
				contentBuilder.WriteString(fmt.Sprintf("Filepath: /%s\n```%s\n%s\n```\n\n", relativePath, filepath.Ext(relativePath)[1:], string(content)))
			} else {
				contentBuilder.WriteString(fmt.Sprintf("Filepath: /%s\n```\n%s\n```\n\n", relativePath, string(content)))
			}
		}

		return nil
	})

	if err != nil {
		return "", "", err
	}

	return structureBuilder.String(), contentBuilder.String(), nil
}

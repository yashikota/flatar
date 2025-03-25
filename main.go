package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer destFile.Close()

	// Copy content
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	err = os.Chmod(dst, srcInfo.Mode())
	if err != nil {
		// Just log permission errors but don't fail
		fmt.Printf("Warning: Could not set permissions for %s: %v\n", dst, err)
	}

	return nil
}

func flattenDirectory(root string) error {
	// Map to keep track of used filenames
	usedNames := make(map[string]bool)

	// First pass: scan for files
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		// Skip the root directory itself
		if path == root {
			return nil
		}

		// Skip directories in this pass
		if info.IsDir() {
			return nil
		}

		// Get relative path from root
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("error getting relative path for %s: %w", path, err)
		}

		// Get filename
		fileName := filepath.Base(path)

		// Check if this filename exists directly in root
		targetPath := filepath.Join(root, fileName)

		// If the file would conflict with an existing file, create a new name
		if _, err := os.Stat(targetPath); err == nil || usedNames[fileName] {
			// If the file is directly in the first level subdirectory, keep the original name
			if filepath.Dir(relPath) == "." {
				// This is a first-level file, keep original name but mark as used
				usedNames[fileName] = true
			} else {
				// Create a unique name by using the directory structure
				ext := filepath.Ext(fileName)
				baseName := fileName[:len(fileName)-len(ext)]

				parentDir := filepath.Base(filepath.Dir(relPath))
				newFileName := fmt.Sprintf("%s_%s%s", baseName, parentDir, ext)

				targetPath = filepath.Join(root, newFileName)
				usedNames[newFileName] = true
			}
		} else {
			usedNames[fileName] = true
		}

		// Copy the file to the new location
		err = copyFile(path, targetPath)
		if err != nil {
			return fmt.Errorf("failed to copy file from %s to %s: %w", path, targetPath, err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Second pass: remove directories (from deepest to shallowest)
	var directories []string

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue despite errors
		}

		if path != root && info.IsDir() {
			directories = append(directories, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Remove directories in reverse order (deepest first)
	for i := len(directories) - 1; i >= 0; i-- {
		err := os.RemoveAll(directories[i])
		if err != nil {
			fmt.Printf("Warning: Failed to remove directory %s: %v\n", directories[i], err)
		}
	}

	return nil
}

func createArchive(sourceDir, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()
	tw := tar.NewWriter(file)
	defer tw.Close()

	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(sourceDir, path)
		hdr := &tar.Header{
			Name: relPath,
			Size: info.Size(),
			Mode: int64(info.Mode()),
		}
		tw.WriteHeader(hdr)
		f, _ := os.Open(path)
		defer f.Close()
		io.Copy(tw, f)
		return nil
	})
}

func main() {
	// Define command line flags
	createArchiveFlag := flag.Bool("a", false, "Create tar archive after flattening")
	deleteDirFlag := flag.Bool("d", false, "Delete original directory after processing")

	// Parse flags
	flag.Parse()

	// Get the root directory from remaining arguments
	args := flag.Args()
	rootDir := "."
	if len(args) == 1 {
		rootDir = args[0]
	} else if len(args) > 1 {
		fmt.Println("Usage: flatar [-a] [-d] [<root_directory>]")
		fmt.Println("  -a    Create tar archive after flattening")
		fmt.Println("  -d    Delete original directory after processing")
		return
	}

	// Ensure the root directory exists and is a directory
	rootInfo, err := os.Stat(rootDir)
	if err != nil {
		fmt.Println("Error accessing directory:", err)
		return
	}
	if !rootInfo.IsDir() {
		fmt.Println("Specified path is not a directory")
		return
	}

	dirs, err := os.ReadDir(rootDir)
	if err != nil {
		fmt.Println("Failed to read root directory:", err)
		return
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		dirPath := filepath.Join(rootDir, dir.Name())
		fmt.Println("Processing:", dirPath)

		if err := flattenDirectory(dirPath); err != nil {
			fmt.Println("Failed to flatten directory:", err)
			continue
		}

		// Create tar archive if -a flag is set
		if *createArchiveFlag {
			tarFile := dirPath + ".tar"
			if err := createArchive(dirPath, tarFile); err != nil {
				fmt.Println("Failed to create archive:", err)
				continue
			}
		}

		// Delete the directory if -d flag is set
		if *deleteDirFlag {
			if err := os.RemoveAll(dirPath); err != nil {
				fmt.Println("Failed to remove directory:", err)
				continue
			}
		}
	}
	fmt.Println("All tasks completed!")
}

package shell

import (
	"context"
	"strings"
)

// handleTabCompletion handles tab completion for files and directories
func (sv *View) handleTabCompletion() {
	currentText := sv.inputField.GetText()
	if currentText == "" {
		return
	}

	words := strings.Fields(currentText)
	if len(words) == 0 {
		return
	}

	if sv.handleSpaceCompletion(currentText) {
		return
	}

	sv.handleWordCompletion(currentText, words)
}

// handleSpaceCompletion handles completion when the input ends with a space
func (sv *View) handleSpaceCompletion(currentText string) bool {
	if !strings.HasSuffix(currentText, " ") {
		return false
	}

	completions := sv.getCompletions(".", "")
	if len(completions) == 0 {
		return true
	}

	if len(completions) == 1 {
		sv.applySingleCompletion(currentText, completions[0])
	} else {
		sv.showCompletions(completions)
	}
	return true
}

// handleWordCompletion handles completion for the last word in the input
func (sv *View) handleWordCompletion(currentText string, words []string) {
	lastWord := words[len(words)-1]
	if lastWord == "" {
		return
	}

	if len(words) == 1 {
		sv.handleCommandCompletionForWords(currentText, lastWord)
		return
	}

	sv.handlePathCompletion(currentText, lastWord)
}

// handleCommandCompletionForWords handles completion for command names
func (sv *View) handleCommandCompletionForWords(currentText, lastWord string) {
	completions := sv.getCommandCompletions(lastWord)
	if len(completions) > 0 {
		sv.handleCommandCompletion(currentText, lastWord, completions)
	}
}

// handlePathCompletion handles completion for file/directory paths
func (sv *View) handlePathCompletion(currentText, lastWord string) {
	dirPath, partialName := sv.parsePathForCompletion(lastWord)
	completions := sv.getCompletions(dirPath, partialName)

	if len(completions) == 0 {
		sv.inputField.SetText(currentText + " ")
		return
	}

	if len(completions) == 1 {
		sv.applySingleCompletion(currentText, completions[0])
	} else {
		sv.handleMultipleCompletions(currentText, partialName, completions)
	}
}

// applySingleCompletion applies a single completion to the input
func (sv *View) applySingleCompletion(currentText, completion string) {
	sv.inputField.SetText(currentText + completion)
	if strings.HasSuffix(completion, "/") {
		sv.inputField.SetText(sv.inputField.GetText() + "/")
	} else {
		sv.inputField.SetText(sv.inputField.GetText() + " ")
	}
}

// handleMultipleCompletions handles multiple completion options
func (sv *View) handleMultipleCompletions(currentText, partialName string, completions []string) {
	commonPrefix := sv.findCommonPrefix(completions)
	if commonPrefix != "" && commonPrefix != partialName {
		newText := strings.TrimSuffix(currentText, partialName) + commonPrefix
		sv.inputField.SetText(newText)
	} else {
		sv.showCompletions(completions)
	}
}

// parsePathForCompletion parses a path to extract directory and partial filename
func (sv *View) parsePathForCompletion(path string) (dir, partial string) {
	if path == "" {
		return ".", ""
	}

	if strings.HasPrefix(path, "/") {
		return sv.parseAbsolutePath(path)
	}

	return sv.parseRelativePath(path)
}

// parseAbsolutePath parses an absolute path
func (sv *View) parseAbsolutePath(path string) (dir, partial string) {
	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == 0 {
		return "/", ""
	}

	dir = path[:lastSlash]
	if dir == "" {
		dir = "/"
	}
	partial = path[lastSlash+1:]
	return dir, partial
}

// parseRelativePath parses a relative path
func (sv *View) parseRelativePath(path string) (dir, partial string) {
	path = strings.TrimPrefix(path, "./")

	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return ".", path
	}

	dir = path[:lastSlash]
	if dir == "" {
		dir = "."
	}
	partial = path[lastSlash+1:]
	return dir, partial
}

// getCompletions queries the container for file/directory completions
func (sv *View) getCompletions(dirPath, partialName string) []string {
	if sv.execFunc == nil {
		return nil
	}

	output, err := sv.executeLsCommand(dirPath)
	if err != nil {
		return nil
	}

	return sv.processLsOutput(output, dirPath, partialName)
}

// executeLsCommand executes the ls command to get directory contents
func (sv *View) executeLsCommand(dirPath string) (string, error) {
	lsArgs := []string{"ls", "-1", dirPath}
	ctx := context.Background()
	return sv.execFunc(ctx, sv.containerID, lsArgs, false)
}

// processLsOutput processes the output of the ls command
func (sv *View) processLsOutput(output, dirPath, partialName string) []string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var completions []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if sv.shouldIncludeCompletion(line, partialName) {
			completion := sv.formatCompletion(line, dirPath)
			completions = append(completions, completion)
		}
	}

	return completions
}

// shouldIncludeCompletion checks if a completion should be included
func (sv *View) shouldIncludeCompletion(line, partialName string) bool {
	return partialName == "" || strings.HasPrefix(line, partialName)
}

// formatCompletion formats a completion line
func (sv *View) formatCompletion(line, dirPath string) string {
	if sv.isDirectory(dirPath, line) {
		return line + "/"
	}
	return line
}

// isDirectory checks if a path is a directory
func (sv *View) isDirectory(dirPath, name string) bool {
	if sv.execFunc == nil {
		return false
	}

	var testPath string
	switch dirPath {
	case ".":
		testPath = name
	case "/":
		testPath = "/" + name
	default:
		testPath = dirPath + "/" + name
	}

	ctx := context.Background()
	_, err := sv.execFunc(ctx, sv.containerID, []string{"ls", "-1", testPath}, false)
	return err == nil
}

// findCommonPrefix finds the common prefix among a list of strings
func (sv *View) findCommonPrefix(strings []string) string {
	if len(strings) == 0 {
		return ""
	}
	if len(strings) == 1 {
		return strings[0]
	}

	minLen := sv.findMinimumLength(strings)
	return sv.buildCommonPrefix(strings, minLen)
}

// findMinimumLength finds the minimum length among all strings
func (sv *View) findMinimumLength(strings []string) int {
	minLen := len(strings[0])
	for _, s := range strings {
		if len(s) < minLen {
			minLen = len(s)
		}
	}
	return minLen
}

// buildCommonPrefix builds the common prefix by comparing characters
func (sv *View) buildCommonPrefix(strings []string, minLen int) string {
	commonPrefix := ""
	for i := 0; i < minLen; i++ {
		if !sv.allStringsHaveSameChar(strings, i) {
			break
		}
		commonPrefix += string(strings[0][i])
	}
	return commonPrefix
}

// allStringsHaveSameChar checks if all strings have the same character at the given index
func (sv *View) allStringsHaveSameChar(strings []string, index int) bool {
	char := strings[0][index]
	for _, s := range strings {
		if s[index] != char {
			return false
		}
	}
	return true
}

// showCompletions displays available completions to the user
func (sv *View) showCompletions(completions []string) {
	sv.addOutput("\nAvailable completions:\n")

	dirs, files := sv.categorizeCompletions(completions)
	sv.displayCompletionsByCategory(dirs, files)
	sv.addOutput("\n")
}

// categorizeCompletions separates completions into directories and files
func (sv *View) categorizeCompletions(completions []string) (dirs, files []string) {
	for _, comp := range completions {
		if strings.HasSuffix(comp, "/") {
			dirs = append(dirs, comp)
		} else {
			files = append(files, comp)
		}
	}
	return dirs, files
}

// displayCompletionsByCategory displays completions grouped by category
func (sv *View) displayCompletionsByCategory(dirs, files []string) {
	sv.displayDirectories(dirs)
	sv.displayFiles(files)
}

// displayDirectories displays directory completions
func (sv *View) displayDirectories(dirs []string) {
	if len(dirs) > 0 {
		sv.addOutput("Directories:\n")
		for _, dir := range dirs {
			sv.addOutput("  " + dir + "\n")
		}
	}
}

// displayFiles displays file completions
func (sv *View) displayFiles(files []string) {
	if len(files) > 0 {
		sv.addOutput("Files:\n")
		for _, file := range files {
			sv.addOutput("  " + file + "\n")
		}
	}
}

// getCommandCompletions gets available command completions
func (sv *View) getCommandCompletions(partial string) []string {
	if sv.execFunc == nil {
		return nil
	}

	commonCommands := sv.getCommonCommands()
	completions := sv.filterCommandsByPartial(commonCommands, partial)
	pathCompletions := sv.getPathCommandCompletions(partial)
	completions = append(completions, pathCompletions...)

	return sv.removeDuplicates(completions)
}

// getCommonCommands returns the list of common commands
func (sv *View) getCommonCommands() []string {
	return []string{
		"ls", "cd", "pwd", "cat", "less", "more", "head", "tail",
		"grep", "find", "which", "whereis", "type", "command",
		"echo", "printf", "date", "whoami", "id", "groups",
		"ps", "top", "htop", "kill", "killall", "pkill",
		"df", "du", "mount", "umount", "fdisk", "blkid",
		"ip", "ifconfig", "netstat", "ss", "ping", "traceroute",
		"curl", "wget", "nc", "telnet", "ssh", "scp",
		"tar", "gzip", "bzip2", "zip", "unzip",
		"vim", "nano", "emacs", "sed", "awk", "sort", "uniq",
		"cut", "paste", "join", "split", "tr", "wc",
		"chmod", "chown", "chgrp", "umask", "touch", "mkdir", "rmdir",
		"cp", "mv", "rm", "ln", "stat", "file", "strings",
		"env", "export", "unset", "set", "alias", "unalias",
		"history", "fc", "jobs", "bg", "fg", "wait",
		"sleep", "timeout", "watch", "nohup", "screen", "tmux",
		"docker", "kubectl", "helm", "git", "svn", "hg",
		"python", "python3", "node", "npm", "java", "javac",
		"gcc", "g++", "make", "cmake", "autoconf", "automake",
		"yum", "apt", "dnf", "pacman", "brew", "snap",
		"systemctl", "service", "init", "systemd", "upstart",
		"cron", "at", "batch", "anacron", "logrotate",
		"rsync", "sftp", "ftp", "tftp", "ncftp",
		"mysql", "psql", "sqlite3", "redis-cli", "mongo",
		"nginx", "apache2ctl", "httpd", "lighttpd", "caddy",
		"fail2ban", "iptables", "ufw", "firewalld", "selinux",
		"auditd", "logwatch", "logcheck", "swatch", "logsurfer",
		"tcpdump", "wireshark", "nmap", "netcat", "socat",
		"strace", "ltrace", "gdb", "valgrind", "perf",
	}
}

// filterCommandsByPartial filters commands that start with the given partial
func (sv *View) filterCommandsByPartial(commands []string, partial string) []string {
	var completions []string
	for _, cmd := range commands {
		if strings.HasPrefix(cmd, partial) {
			completions = append(completions, cmd)
		}
	}
	return completions
}

// getPathCommandCompletions gets commands from PATH
func (sv *View) getPathCommandCompletions(partial string) []string {
	if sv.execFunc == nil {
		return nil
	}

	commonPaths := sv.getCommonPaths()
	return sv.searchPathsForCommands(partial, commonPaths)
}

// getCommonPaths returns the common PATH locations to search
func (sv *View) getCommonPaths() []string {
	return []string{"/usr/bin", "/usr/sbin", "/bin", "/sbin", "/usr/local/bin", "/usr/local/sbin"}
}

// searchPathsForCommands searches the given paths for commands matching the partial
func (sv *View) searchPathsForCommands(partial string, paths []string) []string {
	var completions []string

	for _, path := range paths {
		pathCompletions := sv.searchPathForCommands(partial, path)
		completions = append(completions, pathCompletions...)
	}

	return completions
}

// searchPathForCommands searches a single path for commands matching the partial
func (sv *View) searchPathForCommands(partial, path string) []string {
	ctx := context.Background()
	output, err := sv.execFunc(ctx, sv.containerID, []string{"ls", "-1", path}, false)
	if err != nil {
		return nil
	}

	return sv.processCommandOutput(output, partial)
}

// processCommandOutput processes the command output to find matching completions
func (sv *View) processCommandOutput(output, partial string) []string {
	var completions []string
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasPrefix(line, partial) {
			completions = append(completions, line)
		}
	}

	return completions
}

// handleCommandCompletion handles command name completion
func (sv *View) handleCommandCompletion(currentText, partial string, completions []string) {
	if len(completions) == 0 {
		sv.handleEmptyCompletions(currentText)
		return
	}

	if len(completions) == 1 {
		sv.handleSingleCompletion(currentText, partial, completions[0])
	} else {
		sv.handleMultipleCompletionsForCommands(currentText, partial, completions)
	}
}

// handleEmptyCompletions handles the case when there are no completions
func (sv *View) handleEmptyCompletions(currentText string) {
	sv.inputField.SetText(currentText + " ")
}

// handleSingleCompletion handles the case when there is exactly one completion
func (sv *View) handleSingleCompletion(currentText, partial, completion string) {
	newText := strings.TrimSuffix(currentText, partial) + completion
	sv.inputField.SetText(newText + " ")
}

// handleMultipleCompletionsForCommands handles the case when there are multiple command completions
func (sv *View) handleMultipleCompletionsForCommands(
	currentText, partial string,
	completions []string,
) {
	commonPrefix := sv.findCommonPrefix(completions)
	if commonPrefix != "" && commonPrefix != partial {
		newText := strings.TrimSuffix(currentText, partial) + commonPrefix
		sv.inputField.SetText(newText)
	} else {
		sv.showCommandCompletions(completions)
	}
}

// showCommandCompletions displays available command completions
func (sv *View) showCommandCompletions(completions []string) {
	sv.addOutput("\nAvailable commands:\n")
	for _, cmd := range completions {
		sv.addOutput("  " + cmd + "\n")
	}
	sv.addOutput("\n")
}

// removeDuplicates removes duplicate strings from a slice
func (sv *View) removeDuplicates(strings []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, s := range strings {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}

	return result
}

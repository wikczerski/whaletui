package shell

import (
	"context"
	"strings"
)

// handleTabCompletion handles tab completion for files and directories
func (sv *ShellView) handleTabCompletion() {
	currentText := sv.inputField.GetText()
	if currentText == "" {
		return
	}

	endsWithSpace := strings.HasSuffix(currentText, " ")

	words := strings.Fields(currentText)
	if len(words) == 0 {
		return
	}

	if endsWithSpace {
		completions := sv.getCompletions(".", "")

		if len(completions) == 0 {
			return
		}

		if len(completions) == 1 {
			completion := completions[0]
			sv.inputField.SetText(currentText + completion)

			if strings.HasSuffix(completion, "/") {
				sv.inputField.SetText(sv.inputField.GetText() + "/")
			} else {
				sv.inputField.SetText(sv.inputField.GetText() + " ")
			}
		} else {
			sv.showCompletions(completions)
		}
		return
	}

	lastWord := words[len(words)-1]
	if lastWord == "" {
		return
	}

	if len(words) == 1 {
		completions := sv.getCommandCompletions(lastWord)
		if len(completions) > 0 {
			sv.handleCommandCompletion(currentText, lastWord, completions)
			return
		}
	}

	dirPath, partialName := sv.parsePathForCompletion(lastWord)
	completions := sv.getCompletions(dirPath, partialName)

	if len(completions) == 0 {
		sv.inputField.SetText(currentText + " ")
		return
	}

	if len(completions) == 1 {
		completion := completions[0]
		newText := strings.TrimSuffix(currentText, partialName) + completion
		sv.inputField.SetText(newText)

		if strings.HasSuffix(completion, "/") {
			sv.inputField.SetText(sv.inputField.GetText() + "/")
		} else {
			sv.inputField.SetText(sv.inputField.GetText() + " ")
		}
	} else {
		commonPrefix := sv.findCommonPrefix(completions)
		if commonPrefix != "" && commonPrefix != partialName {
			newText := strings.TrimSuffix(currentText, partialName) + commonPrefix
			sv.inputField.SetText(newText)
		} else {
			sv.showCompletions(completions)
		}
	}
}

// parsePathForCompletion parses a path to extract directory and partial filename
func (sv *ShellView) parsePathForCompletion(path string) (string, string) {
	if path == "" {
		return ".", ""
	}

	if strings.HasPrefix(path, "/") {
		lastSlash := strings.LastIndex(path, "/")
		if lastSlash == 0 {
			return "/", ""
		}
		dir := path[:lastSlash]
		if dir == "" {
			dir = "/"
		}
		partial := path[lastSlash+1:]
		return dir, partial
	}

	path = strings.TrimPrefix(path, "./")

	lastSlash := strings.LastIndex(path, "/")
	if lastSlash == -1 {
		return ".", path
	}

	dir := path[:lastSlash]
	if dir == "" {
		dir = "."
	}
	partial := path[lastSlash+1:]
	return dir, partial
}

// getCompletions queries the container for file/directory completions
func (sv *ShellView) getCompletions(dirPath, partialName string) []string {
	if sv.execFunc == nil {
		return nil
	}

	var lsArgs []string
	if partialName == "" {
		lsArgs = []string{"ls", "-1", dirPath}
	} else {
		lsArgs = []string{"ls", "-1", dirPath}
	}

	ctx := context.Background()
	output, err := sv.execFunc(ctx, sv.containerID, lsArgs, false)
	if err != nil {
		return nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var completions []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if partialName == "" || strings.HasPrefix(line, partialName) {
			if sv.isDirectory(dirPath, line) {
				completions = append(completions, line+"/")
			} else {
				completions = append(completions, line)
			}
		}
	}

	return completions
}

// isDirectory checks if a path is a directory
func (sv *ShellView) isDirectory(dirPath, name string) bool {
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
func (sv *ShellView) findCommonPrefix(strings []string) string {
	if len(strings) == 0 {
		return ""
	}
	if len(strings) == 1 {
		return strings[0]
	}

	minLen := len(strings[0])
	for _, s := range strings {
		if len(s) < minLen {
			minLen = len(s)
		}
	}

	commonPrefix := ""
	for i := 0; i < minLen; i++ {
		char := strings[0][i]
		for _, s := range strings {
			if s[i] != char {
				return commonPrefix
			}
		}
		commonPrefix += string(char)
	}

	return commonPrefix
}

// showCompletions displays available completions to the user
func (sv *ShellView) showCompletions(completions []string) {
	sv.addOutput("\nAvailable completions:\n")

	var dirs, files []string
	for _, comp := range completions {
		if strings.HasSuffix(comp, "/") {
			dirs = append(dirs, comp)
		} else {
			files = append(files, comp)
		}
	}

	if len(dirs) > 0 {
		sv.addOutput("Directories:\n")
		for _, dir := range dirs {
			sv.addOutput("  " + dir + "\n")
		}
	}

	if len(files) > 0 {
		sv.addOutput("Files:\n")
		for _, file := range files {
			sv.addOutput("  " + file + "\n")
		}
	}

	sv.addOutput("\n")
}

// getCommandCompletions gets available command completions
func (sv *ShellView) getCommandCompletions(partial string) []string {
	if sv.execFunc == nil {
		return nil
	}

	// Common commands that are usually available
	commonCommands := []string{
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

	var completions []string
	for _, cmd := range commonCommands {
		if strings.HasPrefix(cmd, partial) {
			completions = append(completions, cmd)
		}
	}

	pathCompletions := sv.getPathCommandCompletions(partial)
	completions = append(completions, pathCompletions...)

	return sv.removeDuplicates(completions)
}

// getPathCommandCompletions gets commands from PATH
func (sv *ShellView) getPathCommandCompletions(partial string) []string {
	if sv.execFunc == nil {
		return nil
	}

	// Try to get commands from common PATH locations
	commonPaths := []string{"/usr/bin", "/usr/sbin", "/bin", "/sbin", "/usr/local/bin", "/usr/local/sbin"}
	var completions []string

	for _, path := range commonPaths {
		ctx := context.Background()
		output, err := sv.execFunc(ctx, sv.containerID, []string{"ls", "-1", path}, false)
		if err != nil {
			continue
		}

		lines := strings.Split(strings.TrimSpace(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && strings.HasPrefix(line, partial) {
				completions = append(completions, line)
			}
		}
	}

	return completions
}

// handleCommandCompletion handles command name completion
func (sv *ShellView) handleCommandCompletion(currentText, partial string, completions []string) {
	if len(completions) == 0 {
		sv.inputField.SetText(currentText + " ")
		return
	}

	if len(completions) == 1 {
		completion := completions[0]
		newText := strings.TrimSuffix(currentText, partial) + completion
		sv.inputField.SetText(newText + " ")
	} else {
		commonPrefix := sv.findCommonPrefix(completions)
		if commonPrefix != "" && commonPrefix != partial {
			newText := strings.TrimSuffix(currentText, partial) + commonPrefix
			sv.inputField.SetText(newText)
		} else {
			sv.showCommandCompletions(completions)
		}
	}
}

// showCommandCompletions displays available command completions
func (sv *ShellView) showCommandCompletions(completions []string) {
	sv.addOutput("\nAvailable commands:\n")
	for _, cmd := range completions {
		sv.addOutput("  " + cmd + "\n")
	}
	sv.addOutput("\n")
}

// removeDuplicates removes duplicate strings from a slice
func (sv *ShellView) removeDuplicates(strings []string) []string {
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

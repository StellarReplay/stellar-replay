package doccheck

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stellar-replay/stellar-replay/internal/cli"
)

func TestMarkdownLinksAndRequiredDocumentationPaths(t *testing.T) {
	root := repositoryRoot(t)
	required := []string{
		"README.md", "CONTRIBUTING.md", "ROADMAP.md", "PROJECT_CONTEXT.md",
		"docs/COMMANDS.md", "docs/DEMO.md", "docs/METHODS.md", "docs/TROUBLESHOOTING.md",
		"docs/FIXTURE_SCHEMA.md", "docs/CAPTURE.md", "docs/REPLAY.md", "docs/SERVER.md",
		"docs/SECURITY.md", "docs/SECURITY_CHECKLIST.md", "examples/README.md",
		"tests/README.md", "fixtures/README.md",
	}
	for _, relative := range required {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(relative))); err != nil {
			t.Errorf("required documentation path missing: %s: %v", relative, err)
		}
	}

	markdownFiles := []string{"README.md", "CONTRIBUTING.md", "ROADMAP.md", "PROJECT_CONTEXT.md"}
	entries, err := os.ReadDir(filepath.Join(root, "docs"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".md") {
			markdownFiles = append(markdownFiles, filepath.ToSlash(filepath.Join("docs", entry.Name())))
		}
	}
	linkPattern := regexp.MustCompile(`\[[^\]]+\]\(([^)]+)\)`)
	for _, relative := range markdownFiles {
		path := filepath.Join(root, filepath.FromSlash(relative))
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range linkPattern.FindAllStringSubmatch(string(data), -1) {
			target := strings.Trim(match[1], "<>")
			if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") || strings.HasPrefix(target, "mailto:") || strings.HasPrefix(target, "#") {
				continue
			}
			target = strings.Split(target, "#")[0]
			if target == "" {
				continue
			}
			resolved := filepath.Clean(filepath.Join(filepath.Dir(path), filepath.FromSlash(target)))
			if _, err := os.Stat(resolved); err != nil {
				t.Errorf("broken documentation link in %s: %s: %v", relative, target, err)
			}
		}
	}
}

func TestDocumentationDoesNotClaimUnimplementedCLI(t *testing.T) {
	root := repositoryRoot(t)
	paths := []string{"README.md", "docs/DEMO.md", "docs/COMMANDS.md", "docs/TROUBLESHOOTING.md"}
	for _, relative := range paths {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			t.Fatal(err)
		}
		text := strings.ToLower(string(data))
		for _, stale := range []string{"commands do not work yet", "exact flags and output are implementation work", "planned commands are"} {
			if strings.Contains(text, stale) {
				t.Errorf("stale implementation claim %q remains in %s", stale, relative)
			}
		}
	}
}

func TestDocumentationListsImplementedCommands(t *testing.T) {
	var output bytes.Buffer
	if code := cli.Run([]string{"help"}, &output, &output); code != cli.ExitOK {
		t.Fatalf("CLI help failed with exit code %d", code)
	}
	for _, command := range []string{"record", "inspect", "validate", "replay", "serve"} {
		if !strings.Contains(output.String(), command) {
			t.Errorf("CLI help omitted implemented command %q", command)
		}
	}
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not locate doccheck source")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", ".."))
}

package agent

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"sync"

	skillsdata "github.com/SparkssL/seer-cli/skills"
)

// discoveredSkill represents a skill found via tree-walking skills.FS.
// The canonical skill ID is the directory name (Dir), not the frontmatter name.
type discoveredSkill struct {
	Dir     string // directory name = canonical skill ID
	Path    string // FS path (e.g., "seer-shared/SKILL.md")
	Name    string // from frontmatter (must match Dir)
	Version string // from frontmatter
}

var (
	discoverOnce sync.Once
	cachedSkills []discoveredSkill
	discoverErr  error
)

// discoverSkills walks skills.FS for */SKILL.md and parses frontmatter.
// Returns an error if discovery fails or any SKILL.md has invalid frontmatter.
func discoverSkills() ([]discoveredSkill, error) {
	discoverOnce.Do(func() {
		matches, err := fs.Glob(skillsdata.FS, "*/SKILL.md")
		if err != nil {
			discoverErr = fmt.Errorf("skill discovery failed: %w", err)
			return
		}
		if len(matches) == 0 {
			discoverErr = fmt.Errorf("no skills found in embedded FS")
			return
		}
		for _, p := range matches {
			data, err := skillsdata.FS.ReadFile(p)
			if err != nil {
				discoverErr = fmt.Errorf("failed to read %s: %w", p, err)
				return
			}
			fm, err := parseFrontmatter(data)
			if err != nil {
				discoverErr = fmt.Errorf("invalid frontmatter in %s: %w", p, err)
				return
			}
			dir := path.Dir(p)
			if fm.Name != dir {
				discoverErr = fmt.Errorf("frontmatter name %q does not match directory %q in %s", fm.Name, dir, p)
				return
			}
			cachedSkills = append(cachedSkills, discoveredSkill{
				Dir:     dir,
				Path:    p,
				Name:    fm.Name,
				Version: fm.Version,
			})
		}
	})
	return cachedSkills, discoverErr
}

// frontmatter holds parsed YAML frontmatter fields from a SKILL.md file.
type frontmatter struct {
	Name        string
	Version     string
	Description string
}

// parseFrontmatter extracts the --- delimited YAML block and parses key fields.
// Only reads top-level scalar fields; nested structures like metadata.requires
// are consumed by skill installers, not by Go code.
func parseFrontmatter(data []byte) (*frontmatter, error) {
	if !bytes.HasPrefix(data, []byte("---\n")) && !bytes.HasPrefix(data, []byte("---\r\n")) {
		return nil, fmt.Errorf("no frontmatter block found")
	}
	rest := data[4:] // skip opening "---\n"
	end := bytes.Index(rest, []byte("\n---"))
	if end < 0 {
		return nil, fmt.Errorf("unclosed frontmatter block")
	}
	block := string(rest[:end])

	fm := &frontmatter{}
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		k, v, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		switch k {
		case "name":
			fm.Name = v
		case "version":
			fm.Version = v
		case "description":
			fm.Description = v
		}
	}
	if fm.Name == "" {
		return nil, fmt.Errorf("frontmatter missing required field: name")
	}
	if fm.Version == "" {
		return nil, fmt.Errorf("frontmatter missing required field: version")
	}
	return fm, nil
}

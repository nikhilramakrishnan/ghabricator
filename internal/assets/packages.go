package assets

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// CelerityMap holds parsed data from resources/celerity/map.php.
type CelerityMap struct {
	// hash → filepath (reverse of 'names' section)
	HashToPath map[string]string
	// symbol → hash
	Symbols map[string]string
	// package → []symbol
	Packages map[string][]string
}

var kvRe = regexp.MustCompile(`'([^']+)'\s*=>\s*'([^']+)'`)
var arrayStartRe = regexp.MustCompile(`'([^']+)'\s*=>\s*array\(`)
var arrayItemRe = regexp.MustCompile(`^\s+'([^']+)',?\s*$`)

// ParseCelerityMap reads and parses resources/celerity/map.php.
func ParseCelerityMap(path string) (*CelerityMap, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read celerity map: %w", err)
	}

	cm := &CelerityMap{
		HashToPath: make(map[string]string),
		Symbols:    make(map[string]string),
		Packages:   make(map[string][]string),
	}

	lines := strings.Split(string(data), "\n")
	section := ""
	var currentPkg string

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Detect top-level section
		if trimmed == "'names' => array(" {
			section = "names"
			continue
		} else if trimmed == "'symbols' => array(" {
			section = "symbols"
			continue
		} else if trimmed == "'requires' => array(" {
			section = "requires"
			continue
		} else if trimmed == "'packages' => array(" {
			section = "packages"
			continue
		}

		// End of a top-level section (only when not inside a nested array)
		if section != "" && currentPkg == "" && trimmed == ")," {
			section = ""
			continue
		}

		switch section {
		case "names":
			if m := kvRe.FindStringSubmatch(line); m != nil {
				path, hash := m[1], m[2]
				cm.HashToPath[hash] = path
			}
		case "symbols":
			if m := kvRe.FindStringSubmatch(line); m != nil {
				sym, hash := m[1], m[2]
				cm.Symbols[sym] = hash
			}
		case "packages":
			if m := arrayStartRe.FindStringSubmatch(line); m != nil {
				currentPkg = m[1]
				cm.Packages[currentPkg] = nil
			} else if currentPkg != "" && trimmed == ")," {
				currentPkg = ""
			} else if currentPkg != "" {
				if m := arrayItemRe.FindStringSubmatch(line); m != nil {
					cm.Packages[currentPkg] = append(cm.Packages[currentPkg], m[1])
				}
			}
		}
	}

	return cm, nil
}

// ResolvePackage returns the ordered list of file paths for a package.
func (cm *CelerityMap) ResolvePackage(pkgName string) ([]string, error) {
	symbols, ok := cm.Packages[pkgName]
	if !ok {
		return nil, fmt.Errorf("unknown package: %s", pkgName)
	}

	var paths []string
	for _, sym := range symbols {
		hash, ok := cm.Symbols[sym]
		if !ok {
			// Skip symbols we can't resolve (e.g. JS-only symbols in CSS context)
			continue
		}
		path, ok := cm.HashToPath[hash]
		if !ok {
			continue
		}
		paths = append(paths, path)
	}
	return paths, nil
}

package gomod

import (
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"os"
	"path"
)

const fileGoMod = "go.mod"

const (
	minSupportedGoVersion = "1.23"    // suffix '.0' omitted on purpose!
	maxSupportedGoVersion = "1.23.99" // allow for future patch versions

	pkgSOM  = "github.com/go-surreal/som"
	pkgSDBC = "github.com/go-surreal/sdbc"

	requiredSOMVersion = "v0.7.1"

	requiredSDBCVersion = "v0.9.2"
)

type GoMod struct {
	path string
	file *modfile.File
}

func NewGoMod(file string, data []byte) (*GoMod, error) {
	modFile, err := modfile.Parse(file, data, nil)
	if err != nil {
		return nil, fmt.Errorf("could not parse go.mod: %v", err)
	}

	return &GoMod{
		path: file,
		file: modFile,
	}, nil
}

func FindGoMod(dir string) (*GoMod, error) {
	var data []byte
	var err error

	for dir != "" {
		data, err = os.ReadFile(path.Join(dir, fileGoMod))

		if err == nil {
			break
		}

		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("could not read go.mod: %v", err)
		}

		dir = path.Dir(dir)
	}

	if data == nil {
		return nil, errors.New("could not find go.mod in worktree")
	}

	return NewGoMod(path.Join(dir, fileGoMod), data)
}

func (m *GoMod) Dir() string {
	return path.Dir(m.path)
}

func (m *GoMod) Module() string {
	return m.file.Module.Mod.Path
}

func (m *GoMod) CheckGoVersion() (string, error) {
	goVersion, err := versionOrdinal(m.file.Go.Version)
	if err != nil {
		return "", fmt.Errorf("could not parse go version: %w", err)
	}

	minSupportedVersion, err := versionOrdinal(minSupportedGoVersion)
	if err != nil {
		return "", fmt.Errorf("could not parse min go version: %w", err)
	}

	maxSupportedVersion, err := versionOrdinal(maxSupportedGoVersion)
	if err != nil {
		return "", fmt.Errorf("could not parse max go version: %w", err)
	}

	if goVersion < minSupportedVersion {
		return "", fmt.Errorf("go version %s is not supported", m.file.Go.Version)
	}

	if goVersion > maxSupportedVersion {
		return fmt.Sprintf("generated code might not work as expected for go version %s (max supported: %s)", m.file.Go.Version, maxSupportedGoVersion), nil
	}

	return "", nil
}

func (m *GoMod) CheckSOMVersion(checkLatest bool) (string, error) {
	for _, require := range m.file.Require {
		if require.Mod.Path != pkgSOM {
			continue
		}

		somVersion, err := versionOrdinal(require.Mod.Version)
		if err != nil {
			return "", fmt.Errorf("could not parse som version: %w", err)
		}

		reqVersion, err := versionOrdinal(requiredSOMVersion)
		if err != nil {
			return "", fmt.Errorf("could not parse required som version: %w", err)
		}

		if somVersion != reqVersion {
			fmt.Printf("go.mod: setting som version to %s\n", requiredSOMVersion)

			if err := m.file.AddRequire(pkgSOM, requiredSOMVersion); err != nil {
				return "", err
			}
		}

		if checkLatest {
			latestVersion, err := SOMVersion()
			if err != nil {
				return "", fmt.Errorf("could not check latest som version: %w", err)
			}

			if somVersion < latestVersion {
				return fmt.Sprintf("newer version of som available: %s (currently: %s)", latestVersion, somVersion), nil
			}
		}

		return "", nil
	}

	fmt.Printf("go.mod: adding som version %s\n", requiredSOMVersion)

	if err := m.file.AddRequire(pkgSOM, requiredSOMVersion); err != nil {
		return "", err
	}

	return "", nil
}

func (m *GoMod) CheckSDBCVersion() (string, error) {
	for _, require := range m.file.Require {
		if require.Mod.Path != pkgSDBC {
			continue
		}

		sdbcVersion, err := versionOrdinal(require.Mod.Version)
		if err != nil {
			return "", fmt.Errorf("could not parse sdbc version: %v", err)
		}

		reqVersion, err := versionOrdinal(requiredSDBCVersion)
		if err != nil {
			return "", fmt.Errorf("could not parse required sdbc version: %v", err)
		}

		if sdbcVersion != reqVersion {
			fmt.Printf("go.mod: setting sdbc version to %s\n", requiredSDBCVersion)

			if err := m.file.AddRequire(pkgSDBC, requiredSDBCVersion); err != nil {
				return "", err
			}

			return "", nil
		}

		return "", nil
	}

	fmt.Printf("go.mod: adding sdbc version %s\n", requiredSDBCVersion)

	if err := m.file.AddRequire(pkgSDBC, requiredSDBCVersion); err != nil {
		return "", err
	}

	return "", nil
}

func (m *GoMod) Save() error {
	content, err := m.file.Format()
	if err != nil {
		return fmt.Errorf("could not format go.mod: %w", err)
	}

	if err := os.WriteFile(m.path, content, 0644); err != nil {
		return err
	}

	return nil
}

//
// -- HELPER
//

func versionOrdinal(version string) (string, error) {
	// ISO/IEC 14651:2011

	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1

	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			return "", fmt.Errorf("invalid version %s", version)
		}
		vo = append(vo, b)
		vo[j]++
	}

	return string(vo), nil
}

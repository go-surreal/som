package util

import (
	"errors"
	"fmt"
	"golang.org/x/mod/modfile"
	"os"
	"path"
)

const fileGoMod = "go.mod"

const (
	minUnsupportedGoVersion = "1.20.0"
	minSupportedGoVersion   = "1.21.9"
	maxSupportedGoVersion   = "1.22.99" // allow for future patch versions

	pkgSOM  = "github.com/go-surreal/som"
	pkgSDBC = "github.com/go-surreal/sdbc"

	minSupportedSOMVersion = "v0.4.0"
	maxSupportedSOMVersion = "v0.4.0"

	minSupportedSDBCVersion = "v0.3.0"
	maxSupportedSDBCVersion = "v0.3.0"
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
		return "", fmt.Errorf("could not parse go version: %v", err)
	}

	minUnsupportedVersion, err := versionOrdinal(minUnsupportedGoVersion)
	if err != nil {
		return "", fmt.Errorf("could not parse min go version: %v", err)
	}

	minSupportedVersion, err := versionOrdinal(minSupportedGoVersion)
	if err != nil {
		return "", fmt.Errorf("could not parse min go version: %v", err)
	}

	maxSupportedVersion, err := versionOrdinal(maxSupportedGoVersion)
	if err != nil {
		return "", fmt.Errorf("could not parse max go version: %v", err)
	}

	if goVersion < minUnsupportedVersion {
		return "", fmt.Errorf("go version %s is not supported", goVersion)
	}

	if goVersion < minSupportedVersion || goVersion > maxSupportedVersion {
		return fmt.Sprintf("generated code might not work as expected for go version %s", goVersion), nil
	}

	return "", nil
}

func (m *GoMod) CheckSOMVersion() (string, error) {
	for _, require := range m.file.Require {
		if require.Mod.Path != pkgSOM {
			continue
		}

		somVersion, err := versionOrdinal(require.Mod.Version)
		if err != nil {
			return "", fmt.Errorf("could not parse som version: %w", err)
		}

		minVersion, err := versionOrdinal(minSupportedSOMVersion)
		if err != nil {
			return "", fmt.Errorf("could not parse min som version: %w", err)
		}

		maxVersion, err := versionOrdinal(maxSupportedSOMVersion)
		if err != nil {
			return "", fmt.Errorf("could not parse max som version: %w", err)
		}

		if somVersion < minVersion {
			return "", fmt.Errorf("som version %s is not supported", require.Mod.Version)
		}

		if somVersion > maxVersion {
			return fmt.Sprintf("generated code might not work as expected for som version %s", require.Mod.Version), nil
		}

		latestVersion, err := SOMVersion()
		if err != nil {
			return "", fmt.Errorf("could not check latest som version: %w", err)
		}

		if somVersion < latestVersion {
			fmt.Println(somVersion, latestVersion)
			return fmt.Sprintf("newer version of som available: %s", latestVersion), nil
		}

		return "", nil
	}

	return "", errors.New("could not find som package in go.mod")
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

		minVersion, err := versionOrdinal(minSupportedSDBCVersion)
		if err != nil {
			return "", fmt.Errorf("could not parse min sdbc version: %v", err)
		}

		maxVersion, err := versionOrdinal(maxSupportedSDBCVersion)
		if err != nil {
			return "", fmt.Errorf("could not parse max sdbc version: %v", err)
		}

		if sdbcVersion < minVersion {
			return "", fmt.Errorf("sdbc version %s is not supported", require.Mod.Version)
		}

		if sdbcVersion > maxVersion {
			return fmt.Sprintf("generated code might not work as expected for sdbc version %s", require.Mod.Version), nil
		}

		latestVersion, err := SDBCVersion()
		if err != nil {
			return "", fmt.Errorf("could not check latest sdbc version: %w", err)
		}

		if sdbcVersion < latestVersion {
			return fmt.Sprintf("newer version of sdbc available: %s", latestVersion), nil
		}

		return "", nil
	}

	return "", errors.New("could not find sdbc package in go.mod")
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

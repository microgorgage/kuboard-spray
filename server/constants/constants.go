package constants

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

func GET_DATA_CLUSTER_DIR() string {
	return GET_DATA_DIR() + "/cluster"
}

func GET_DATA_RESOURCE_DIR() string {
	return GET_DATA_DIR() + "/resource"
}

func GET_DATA_MIRROR_DIR() string {
	return GET_DATA_DIR() + "/mirror"
}

func GET_DATA_DIR() string {
	val, ex := os.LookupEnv("KUBOARD_SPRAY_DATA_DIR")
	if !ex {
		home, err := Home()
		if err != nil {
			home = "~/"
		}
		return home + "/.kuboardspray-data"
	}
	return val
}

func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support
	if runtime.GOOS == "windows" {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

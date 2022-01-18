package pathUtils

import (
	"crypto-bot/pkg/utils/envUtils"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"
)

// GetRootProjectPath return root path of project (ex: /home/$USER/Documents/starter-module)
// APPLICATION_NAME env variable must be set in order to get root project directory (ex: starter-module is the name of git project)
// INIT_ROOT_PROJECT_PATH can be used to override root project path
func GetRootProjectPath() (rootPath string, err error) {
	rootPath = envUtils.GetFromEnvOrDefault(envInitRootProjectPath, defaultInitRootProjectPath)
	if rootPath != defaultInitRootProjectPath {
		return
	}

	projectName := "crypto-bot"

	// Get current path
	wd, err := os.Getwd()
	if err != nil {
		return rootPath, errors.WithStack(err)
	}

	if wd == "/" {
		rootPath = "/"
		return
	}

	// Find root project folder based on project name
	rootDirIndex := strings.Index(wd, projectName)
	if rootDirIndex < 0 {
		return rootPath, fmt.Errorf("cannot get root project path from wd %s", wd)
	}
	rootDir := wd[:rootDirIndex]
	rootPath = rootDir + projectName

	return
}

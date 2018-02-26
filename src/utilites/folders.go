package utilites

import (
  "os/exec"
  "color"
  "fmt"
  "os"
	"path/filepath"
  "strings"
)

func GetDownloadsDirectory() string {
  out, err := exec.Command("sh", "-c", "xdg-user-dir DOWNLOAD").Output()

  var path string
  if err != nil {
    fmt.Println(color.Red("Error getting downloads directory path :/  "), err)
  } else {
    path = string(out)
  }

  return strings.TrimSpace(path)
}

func GetWorkingDirectory() string {
  wd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
  return wd
}

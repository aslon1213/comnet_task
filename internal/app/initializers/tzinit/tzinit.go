package tzinit

import "os"

func init() {
	os.Setenv("TZ", "Asia/Tashkent")
}

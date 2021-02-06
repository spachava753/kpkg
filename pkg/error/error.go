package error

import (
	"fmt"
	"runtime"
)

type UnsupportedRuntimeErr struct {
	Binary string
}

func (u *UnsupportedRuntimeErr) Error() string {
	return fmt.Sprintf("downloading binary %s is not support with runtime %s/%s", u.Binary, runtime.GOOS, runtime.GOARCH)
}

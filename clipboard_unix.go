// +build freebsd linux netbsd openbsd solaris dragonfly

package main

import "github.com/atotto/clipboard"

var hasPrimary = true

func setPrimary(enabled bool) {
	clipboard.Primary = enabled

}

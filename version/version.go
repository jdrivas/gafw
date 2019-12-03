package version

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// var gitHash = "EMPTY"
// var gitTag = "v0.0.0"
// var unixtime = "EMPTY"

var gitHash string
var gitTag string
var unixtime string

// Version is used to obtain compile time information.
var Version AppVersion

// AppVersion caputres compile time infromation for display
type AppVersion struct {
	Major, Minor, Dot int
	Hash              string
	Date              time.Time
}

func init() {
	mj, mi, d := tagToSem(gitTag)
	ut, _ := strconv.ParseInt(unixtime, 10, 64)
	bt := time.Unix(ut, 0)
	Version = AppVersion{mj, mi, d, gitHash, bt}
}

// Expect gt => v{int}.{int}.{int}
func tagToSem(t string) (mj, mi, d int) {
	if t == "" {
		return 0, 0, 0
	}
	ts := strings.SplitN(t[1:], ".", 3) // point past the initial. "v"
	var vals [3]int
	for i, v := range ts {
		vals[i], _ = strconv.Atoi(v)
	}
	return vals[0], vals[1], vals[2]
}

// String pretty prints version information
func (v AppVersion) String() string {
	var s string
	if v.Major == 0 && v.Minor == 0 && v.Dot == 0 {
		s = fmt.Sprintf("Version: [unreleased] [%s] %s", v.Hash, v.Date.Local().Format(time.RFC1123))
	} else {
		s = fmt.Sprintf("Version: %d.%d.%d [%s] %s", v.Major, v.Minor, v.Dot, v.Hash, v.Date.Local().Format(time.RFC1123))
	}
	return s
}

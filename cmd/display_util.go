package cmd

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"

	"github.com/jdrivas/conman"
	t "github.com/jdrivas/termtext"
	"github.com/juju/ansiterm"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Shim for t.httpDisplay, with timing display.
func httpDisplay(se *conman.SideEffect, resp *http.Response, err error) {
	if !viper.GetBool(t.JSONDisplayKey) {
		if se.ElapsedTime.Milliseconds() < 1000 {
			fmt.Printf(t.Title("Command took %d milliseconds\n", se.ElapsedTime.Milliseconds()))
		} else {
			fmt.Printf(t.Title("Command took %4g seconds\n", se.ElapsedTime.Seconds()))

		}
	}
	t.HTTPDisplay(resp, err)
}

// Configuration

func printConfig() {
	settings := viper.AllSettings()
	fmt.Printf("%s\n", t.SubTitle("Settings are:"))
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 2, ' ', 0)
	fmt.Fprintf(w, "%s", t.Title("Name\tValue\n"))
	keys := []string{}
	for k := range settings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, configString(k, settings[k], 1))
	}
	w.Flush()
}

const maxlen = 60

func configString(k string, v interface{}, depth int) (rs string) {
	// fmt.Printf("configString: %s(%d): %v\n", k, depth, v)
	ts := ""
	for i := 1; i < depth; i++ {
		ts += "\t"
	}

	switch sv := v.(type) {
	case bool:
		rs += ts + t.Title("%s:\t", k) + t.SubTitle("%t\n", sv)
	case int:
		rs += ts + t.Title("%s:\t", k) + t.SubTitle("%d\n", sv)
	case string:
		if len(sv) > maxlen {
			sv = sv[:maxlen/2] + "..." + sv[len(sv)-maxlen/2:]
		}
		rs += ts + t.Title("%s:\t", k) + t.SubTitle("%s\n", sv)
	case map[string]interface{}:
		rs += ts + t.Title("%s:\t\n", k)
		keys := []string{}
		for k1 := range sv {
			keys = append(keys, k1)
		}
		sort.Strings(keys)
		for _, k1 := range keys {
			rs += configString(k1, sv[k1], depth+1)
		}
	default:
		rs += ts + t.Title("%s:\t", k) + t.SubTitle("%#v\n", sv)

	}
	return rs
}

// Flags
/*
func (sfm bindFlagMap) print() {

	keys := []string{}
	for k := range sfm {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 2, ' ', 0)
	fmt.Fprintf(w, flagHeader()+t.Title("\tSaved-Value\tBindKey\tConvert\n")) // Add the saveFlag values
	for _, k := range keys {
		f := sfm[k]
		n := nameOf(f.Convert)
		fmt.Fprintf(w, flagEntry(f.Flag)+t.SubTitle("\t%v\t%s\t%s\n", f.Value, f.BindKey, n))
	}
	w.Flush()
}
*/

func printFlagSet(flags *pflag.FlagSet) {
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 2, ' ', 0)
	fmt.Fprintf(w, flagHeader()+"\n")
	flags.VisitAll(func(f *pflag.Flag) {
		fmt.Fprintf(w, flagEntry(f)+"\n")
	})
	w.Flush()
}

func flagHeader() string {
	return t.Title("Name\tShort\tValue\tType\tDefValue\tChanged")
}

func flagEntry(f *pflag.Flag) string {
	if f != nil {
		return t.SubTitle("%s\t%s\t%s\t%s\t%s\t%t",
			f.Name, f.Shorthand, f.Value.String(), f.Value.Type(), f.DefValue, f.Changed)
	}
	return "<No Flag>\t-\t-\t-\t-\t-"
}

//
func nameOf(f interface{}) string {
	v := reflect.ValueOf(f)
	return filepath.Base(runtime.FuncForPC(v.Pointer()).Name())
}

// Stack tracing
// look at runtime/debug PrintStack() for a full stack trace dump.

// print info about where the call is.
// depth = 0 is the current location.
// depth = 1 is one caller up etc.
func printLoc(mesg string, depth int) {
	fc, fl, ln := locString(depth + 1) // print from the caller.
	fmt.Printf("%s\t%s()\t%s:%d\n", mesg, fc, fl, ln)
}

func locString(d int) (fnc, file string, line int) {
	if pc, fl, l, ok := runtime.Caller(d + 1); ok { // add one, because 0 should be the coller of this funciton.
		f := runtime.FuncForPC(pc)
		fnc = filepath.Base(f.Name())
		file = filepath.Base(fl)
		line = l
	}
	return fnc, file, line
}

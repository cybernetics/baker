package testutil

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

// TempDir is a test helper that creates a temporary directory, returns its
// name and a function which when called removes that directory. This is useful
// to be called as follows at the top of a test or benchmark requiring a
// temporary directory:
//
//  func TestFoo(t *testing.T) {
//      dir, rmdir := testutil.TempDir(t)
//      defer rmdir()
//
//      // do something with dir
//      ...
//  }
func TempDir(tb testing.TB) (dir string, rmdir func()) {
	tb.Helper()

	var err error
	if dir, err = ioutil.TempDir("", tb.Name()); err != nil {
		tb.Fatalf("can't create temp directory: %v", err)
	}

	rmdir = func() {
		if err = os.RemoveAll(dir); err != nil {
			tb.Fatalf("can't remove temp directory: %v", err)
		}
	}
	return dir, rmdir
}

// TempFile is a test helper that creates a temporary file, returns its
// name and a function which when called removes that file. This is useful
// to be called as follows at the top of a test or benchmark requiring a
// temporary file:
//
//  func TestFoo(t *testing.T) {
//      name, rmfile := testutil.TempFile(t)
//      defer rmfile()
//
//      // do something with filename
//      ...
//  }
func TempFile(tb testing.TB) (file string, rmfile func()) {
	tb.Helper()

	f, err := ioutil.TempFile("", tb.Name())
	if err != nil {
		tb.Fatalf("can't create temp file: %v", err)
	}

	if err = f.Close(); err != nil {
		tb.Fatalf("can't create temp file: %v", err)
	}

	rmfile = func() {
		if err = os.Remove(f.Name()); err != nil {
			tb.Fatalf("can't remove temp file: %v", err)
		}
	}
	return f.Name(), rmfile
}

// DisableLogging is a test helper that disable logging (in fact it sets its
// level to panic). It returns a function which when called, resets it to its
// previous level. Its useful to be called as follows in test/benchmarks:
//
//  func TestFoo(t *testing.T) {
//      defer DisableLogging()()
//
//      // logging is disabled for the whole test
//  }
func DisableLogging() (reset func()) {
	lvl := logrus.GetLevel()
	logrus.SetLevel(logrus.PanicLevel)
	return func() { logrus.SetLevel(lvl) }
}

// LessLogging is a test helper that decreases logging (in fact it sets its
// level to Error). It returns a function which when called, resets it to its
// previous level. Its useful to be called as follows in test/benchmarks:
//
//  func TestFoo(t *testing.T) {
//      defer LessLogging()()
//
//      // logging is set to Error for the whole test
//  }
func LessLogging() (reset func()) {
	lvl := logrus.GetLevel()
	logrus.SetLevel(logrus.ErrorLevel)
	return func() { logrus.SetLevel(lvl) }
}

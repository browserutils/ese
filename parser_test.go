package ese

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	update = flag.Bool("update", false, "update golden files")
	binary string
	assert assertRunner
	goldie goldieRunner
	suite  suiteRunner
)

type testSuite struct {
	t *testing.T
}

func (self *testSuite) T() *testing.T {
	return self.t
}

type ESETestSuite struct {
	testSuite
	binary string
}

type assertRunner struct{}

func (assertRunner) NoError(t *testing.T, err error, args ...interface{}) {
	if err == nil {
		return
	}

	if len(args) > 0 {
		t.Fatal(args...)
	}
	t.Fatal(err)
}

type goldieRunner struct{}

func (goldieRunner) Assert(t *testing.T, fixture_name string, out []byte) {
	t.Helper()

	golden := filepath.Join("testdata", "fixtures", fixture_name+".golden")
	if *update {
		err := os.WriteFile(golden, out, 0644)
		if err != nil {
			t.Fatal(err)
		}
		return
	}

	expected, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, out) {
		t.Fatalf("golden mismatch for %s", fixture_name)
	}
}

type suiteRunner struct{}

func TestMain(m *testing.M) {
	binary = "./eseparser"
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binary, "./cmd/eseparser")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(out))
		os.Exit(1)
	}

	code := m.Run()
	_ = os.Remove(binary)
	os.Exit(code)
}

func (self *ESETestSuite) SetupTest() {
	self.binary = "./eseparser"
	if runtime.GOOS == "windows" {
		self.binary += ".exe"
	}
}

func (suiteRunner) Run(t *testing.T, self *ESETestSuite) {
	t.Run("UAL", func(t *testing.T) {
		self.t = t
		self.SetupTest()
		self.TestUAL()
	})
	t.Run("SystemIdentity", func(t *testing.T) {
		self.t = t
		self.SetupTest()
		self.TestSystemIdentity()
	})
	t.Run("SRUM", func(t *testing.T) {
		self.t = t
		self.SetupTest()
		self.TestSRUM()
	})
	t.Run("Ntds", func(t *testing.T) {
		self.t = t
		self.SetupTest()
		self.TestNtds()
	})
	t.Run("NtdsLongValues", func(t *testing.T) {
		self.t = t
		self.SetupTest()
		self.TestNtdsLongValues()
	})
	t.Run("WebCache", func(t *testing.T) {
		self.t = t
		self.SetupTest()
		self.TestWebCache()
	})
	t.Run("WindowsEdb", func(t *testing.T) {
		self.t = t
		self.SetupTest()
		self.TestWindowsEdb()
	})
	t.Run("WindowsQmgr", func(t *testing.T) {
		self.t = t
		self.SetupTest()
		self.TestWindowsQmgr()
	})
}

// User Access Logs have some interesting columns types:
//   - GUID
//   - DateTime seem to be encoded in a different way - a uint64 windows
//     file time.
func (self *ESETestSuite) TestUAL() {
	cmdline := []string{
		"dump", "--limit", "2",
		"testdata/Sample_UAL/HyperV-PC/Current.mdb", "CLIENTS",
	}
	cmd := exec.Command(self.binary, cmdline...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
	}
	assert.NoError(self.T(), err)

	fixture_name := "UAL_CLIENTS"
	goldie.Assert(self.T(), fixture_name, out)
}

func (self *ESETestSuite) TestSystemIdentity() {
	cmdline := []string{
		"dump",
		"./testdata/Sample_UAL/HyperV-PC/SystemIdentity.mdb", "SYSTEM_IDENTITY",
	}
	cmd := exec.Command(self.binary, cmdline...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
	}
	assert.NoError(self.T(), err)

	fixture_name := "SYSTEM_IDENTITY"
	goldie.Assert(self.T(), fixture_name, out)
}

func (self *ESETestSuite) TestSRUM() {
	cmdline := []string{
		"dump", "--limit", "2",
		"testdata/SRUM/SRUDB.dat", "{D10CA2FE-6FCF-4F6D-848E-B2E99266FA86}",
	}
	cmd := exec.Command(self.binary, cmdline...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		assert.NoError(self.T(), err)
	}
	fixture_name := "SRUM-D10CA2FE-6FCF-4F6D-848E-B2E99266FA86"
	goldie.Assert(self.T(), fixture_name, out)
}

func (self *ESETestSuite) TestNtds() {
	cmdline := []string{
		"dump", "--limit", "5",
		"testdata/Samples/ntds.dit", "datatable",
	}
	cmd := exec.Command(self.binary, cmdline...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		assert.NoError(self.T(), err)
	}
	fixture_name := "ntds.dit"
	goldie.Assert(self.T(), fixture_name, out)
}

func (self *ESETestSuite) TestNtdsLongValues() {
	cmdline := []string{
		"dump", "--limit", "500",
		"testdata/Samples/ntds.dit", "sd_table",
	}
	cmd := exec.Command(self.binary, cmdline...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		assert.NoError(self.T(), err)
	}
	fixture_name := "ntds.dit_sd_table"
	goldie.Assert(self.T(), fixture_name, out)
}

func (self *ESETestSuite) TestWebCache() {
	cmdline := []string{
		"dump", "testdata/Samples/WebCacheV01.dat", "Containers", "Container_2",
	}
	cmd := exec.Command(self.binary, cmdline...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		assert.NoError(self.T(), err, string(out))
	}
	fixture_name := "WebCacheV01.dat"
	goldie.Assert(self.T(), fixture_name, out)
}

func (self *ESETestSuite) TestWindowsEdb() {
	cmdline := []string{
		"dump", "testdata/Samples/Windows.edb",
		"SystemIndex_Gthr", "SystemIndex_GthrPth", "SystemIndex_PropertyStore",
		"--limit", "10",
	}
	cmd := exec.Command(self.binary, cmdline...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		assert.NoError(self.T(), err, string(out))
	}
	fixture_name := "WindowsEdb"
	goldie.Assert(self.T(), fixture_name, out)
}

func (self *ESETestSuite) TestWindowsQmgr() {
	cmdline := []string{
		"dump", "testdata/Samples/qmgr.db",
		"--limit", "10",
	}
	cmd := exec.Command(self.binary, cmdline...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		assert.NoError(self.T(), err, string(out))
	}
	fixture_name := "WindowsQmgr"
	goldie.Assert(self.T(), fixture_name, out)
}

func TestESE(t *testing.T) {
	suite.Run(t, &ESETestSuite{})
}

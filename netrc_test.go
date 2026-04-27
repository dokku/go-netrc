package netrc_test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jdx/go-netrc"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type NetrcSuite struct{}

var _ = Suite(&NetrcSuite{})

func (s *NetrcSuite) TestLogin(c *C) {
	f, err := netrc.Parse("./examples/login.netrc")
	c.Assert(err, IsNil)
	heroku := f.Machine("api.heroku.com")
	c.Check(heroku.Get("login"), Equals, "jeff@heroku.com")
	c.Check(heroku.Get("password"), Equals, "foo")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestSave(c *C) {
	f, err := netrc.Parse("./examples/login.netrc")
	c.Assert(err, IsNil)
	f.Path = "./examples/login-new.netrc"
	err = f.Save()
	c.Assert(err, IsNil)
	a, _ := ioutil.ReadFile("./examples/login-new.netrc")
	b, _ := ioutil.ReadFile("./examples/login.netrc")
	c.Check(string(a), Equals, string(b))
	os.Remove("./examples/login-new.netrc")
}

func (s *NetrcSuite) TestAdd(c *C) {
	f, err := netrc.Parse("./examples/login.netrc")
	c.Assert(err, IsNil)
	f.AddMachine("m", "l", "p")
	c.Check(f.Render(), Equals, "# this is my login netrc\nmachine api.heroku.com\n  login jeff@heroku.com # this is my username\n  password foo\n"+
		"machine m\n  login l\n  password p\n")
}

func (s *NetrcSuite) TestMachineSetWithEmptyFile(c *C) {
	dir := c.MkDir()
	n := netrc.New(filepath.Join(dir, ".netrc"))
	n.AddMachine("api.heroku.com", "jeff@heroku.com", "foo")
	c.Check(n.Render(), Equals, "machine api.heroku.com\n  login jeff@heroku.com\n  password foo\n")
}

func (s *NetrcSuite) TestAddExisting(c *C) {
	f, err := netrc.Parse("./examples/login.netrc")
	c.Assert(err, IsNil)
	f.AddMachine("api.heroku.com", "l", "p")
	c.Check(f.Render(), Equals, "# this is my login netrc\nmachine api.heroku.com\n  login l\n  password p\n")
}

func (s *NetrcSuite) TestRemove(c *C) {
	f, err := netrc.Parse("./examples/sample_multi.netrc")
	c.Assert(err, IsNil)
	f.RemoveMachine("m")
	c.Check(f.Render(), Equals, "# this is my netrc with multiple machines\nmachine n\n  login ln # this is my n-username\n  password pn\n")
}

func (s *NetrcSuite) TestSetPassword(c *C) {
	f, err := netrc.Parse("./examples/login.netrc")
	c.Assert(err, IsNil)
	heroku := f.Machine("api.heroku.com")
	heroku.Set("password", "foobar")
	c.Check(f.Render(), Equals, "# this is my login netrc\nmachine api.heroku.com\n  login jeff@heroku.com # this is my username\n  password foobar\n")
}

func (s *NetrcSuite) TestSampleMulti(c *C) {
	f, err := netrc.Parse("./examples/sample_multi.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("m").Get("login"), Equals, "lm")
	c.Check(f.Machine("m").Get("password"), Equals, "pm")
	c.Check(f.Machine("n").Get("login"), Equals, "ln")
	c.Check(f.Machine("n").Get("password"), Equals, "pn")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestSampleMultiWithDefault(c *C) {
	f, err := netrc.Parse("./examples/sample_multi_with_default.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("m").Get("login"), Equals, "lm")
	c.Check(f.Machine("m").Get("password"), Equals, "pm")
	c.Check(f.Machine("n").Get("login"), Equals, "ln")
	c.Check(f.Machine("n").Get("password"), Equals, "pn")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestNewlineless(c *C) {
	f, err := netrc.Parse("./examples/newlineless.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("m").Get("login"), Equals, "l")
	c.Check(f.Machine("m").Get("password"), Equals, "p")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestBadDefaultOrder(c *C) {
	f, err := netrc.Parse("./examples/bad_default_order.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("mail.google.com").Get("login"), Equals, "joe@gmail.com")
	c.Check(f.Machine("mail.google.com").Get("password"), Equals, "somethingSecret")
	c.Check(f.Machine("ray").Get("login"), Equals, "demo")
	c.Check(f.Machine("ray").Get("password"), Equals, "mypassword")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestDefaultOnly(c *C) {
	f, err := netrc.Parse("./examples/default_only.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("default").Get("login"), Equals, "ld")
	c.Check(f.Machine("default").Get("password"), Equals, "pd")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestGood(c *C) {
	f, err := netrc.Parse("./examples/good.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("mail.google.com").Get("login"), Equals, "joe@gmail.com")
	c.Check(f.Machine("mail.google.com").Get("account"), Equals, "justagmail")
	c.Check(f.Machine("mail.google.com").Get("password"), Equals, "somethingSecret")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestPassword(c *C) {
	f, err := netrc.Parse("./examples/password.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("m").Get("password"), Equals, "p")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestGetOctothorpe(c *C) {
	f, err := netrc.Parse("./examples/octothorpe.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("hash").Get("password"), Equals, "foo#bar$baz%boom")
	c.Check(f.Machine("hash2").Get("password"), Equals, "foo#bar$baz%boom##")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestPermissive(c *C) {
	f, err := netrc.Parse("./examples/permissive.netrc")
	c.Assert(err, IsNil)
	c.Check(f.Machine("m").Get("login"), Equals, "l")
	c.Check(f.Machine("m").Get("password"), Equals, "p")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestPasswordWithSpaces(c *C) {
	f, err := netrc.Parse("./examples/spaces.netrc")
	c.Assert(err, IsNil)
	m := f.Machine("spaces.example.com")
	c.Check(m.Get("login"), Equals, "alice")
	c.Check(m.Get("password"), Equals, "my pass with spaces")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestPasswordWithEscapes(c *C) {
	f, err := netrc.Parse("./examples/escaped.netrc")
	c.Assert(err, IsNil)
	m := f.Machine("escaped.example.com")
	c.Check(m.Get("login"), Equals, "alice")
	c.Check(m.Get("password"), Equals, "say \"hi\"\nthere")
	body, _ := ioutil.ReadFile(f.Path)
	c.Check(f.Render(), Equals, string(body))
}

func (s *NetrcSuite) TestSetPasswordWithSpaces(c *C) {
	f, err := netrc.Parse("./examples/login.netrc")
	c.Assert(err, IsNil)
	heroku := f.Machine("api.heroku.com")
	heroku.Set("password", "my new pass")
	c.Check(heroku.Get("password"), Equals, "my new pass")
	c.Check(f.Render(), Equals, "# this is my login netrc\nmachine api.heroku.com\n  login jeff@heroku.com # this is my username\n  password \"my new pass\"\n")
}

func (s *NetrcSuite) TestSetPasswordWithQuoteAndNewline(c *C) {
	dir := c.MkDir()
	n := netrc.New(filepath.Join(dir, ".netrc"))
	n.AddMachine("m", "alice", "ab\"c\nd")
	rendered := n.Render()
	c.Check(rendered, Equals, "machine m\n  login alice\n  password \"ab\\\"c\\nd\"\n")
	round, err := netrc.ParseString(rendered)
	c.Assert(err, IsNil)
	c.Check(round.Machine("m").Get("password"), Equals, "ab\"c\nd")
}

func (s *NetrcSuite) TestSetPasswordNoQuoteWhenSafe(c *C) {
	f, err := netrc.Parse("./examples/login.netrc")
	c.Assert(err, IsNil)
	heroku := f.Machine("api.heroku.com")
	heroku.Set("password", "plain")
	c.Check(f.Render(), Equals, "# this is my login netrc\nmachine api.heroku.com\n  login jeff@heroku.com # this is my username\n  password plain\n")
}

func (s *NetrcSuite) TestAddMachineWithSpaces(c *C) {
	dir := c.MkDir()
	n := netrc.New(filepath.Join(dir, ".netrc"))
	n.AddMachine("m", "user", "pass with spaces")
	round, err := netrc.ParseString(n.Render())
	c.Assert(err, IsNil)
	c.Check(round.Machine("m").Get("login"), Equals, "user")
	c.Check(round.Machine("m").Get("password"), Equals, "pass with spaces")
}

func (s *NetrcSuite) TestParseString(c *C) {
	file, err := os.Open("./examples/good.netrc")
	defer file.Close()
	data, err := io.ReadAll(file)
	c.Assert(err, IsNil)
	f, err := netrc.ParseString(string(data))
	c.Assert(err, IsNil)
	c.Check(f.Machine("mail.google.com").Get("login"), Equals, "joe@gmail.com")
	c.Check(f.Machine("mail.google.com").Get("account"), Equals, "justagmail")
	c.Check(f.Machine("mail.google.com").Get("password"), Equals, "somethingSecret")
}

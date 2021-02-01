package pdfParse

import (
	"os"
	"path"
	"reflect"
	"testing"
)

var contentdir string = "ContentTest"
var fontfilename string = "saxmono.ttf"
var footerfilename string = "logo.jpg"
var inputfilename string = "0433 Race  Need4Speed.pdf"
var outfilename string = "out.jpg"

func TestNewRaceSession(t *testing.T) {
	var rs *RaceSession
	var err error
	params := Parameters{}
	rs, err = NewRaceSession(params)
	if err == nil {
		t.Logf("%+v", rs)
		t.Fail()
	}
	params.InputfilePath = "test/path"
	rs, err = NewRaceSession(params)
	if err == nil {
		t.Logf("%+v", rs)
		t.Fail()
	}
	params.FontfilePath = "test/path"
	rs, err = NewRaceSession(params)
	if err == nil {
		t.Logf("%+v", rs)
		t.Fail()
	}
	params.OutputimagePath = "test/path"
	rs, err = NewRaceSession(params)
	if err != nil {
		t.Error(err)
	}
	if reflect.TypeOf(rs) != reflect.TypeOf(&RaceSession{}) {
		t.Error(reflect.TypeOf(rs))
	}
}

func TestPdfToImage(t *testing.T) {
	inputfilename = "0433 Race  Need4Speed.pdf"
	outfilename = "out.jpg"
	outpath := path.Join(contentdir, outfilename)
	params := Parameters{InputfilePath: path.Join(contentdir, inputfilename),
		FontfilePath:    path.Join(contentdir, fontfilename),
		FooterImagePath: path.Join(contentdir, footerfilename),
		OutputimagePath: outpath,
		FontSize:        5.5,
	}

	os.Remove(outpath)

	rs, err := NewRaceSession(params)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	err = rs.PdfToImage()
	if err != nil {
		t.Error(err)
	}
}

func TestSessionToLines(t *testing.T) {
	outpath := path.Join(contentdir, outfilename)
	params := Parameters{InputfilePath: path.Join(contentdir, inputfilename),
		FontfilePath:    path.Join(contentdir, fontfilename),
		FooterImagePath: path.Join(contentdir, footerfilename),
		OutputimagePath: outpath}
	rs, _ := NewRaceSession(params)
	err := rs.ReadPdf()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	var lines []string

	_, err = rs.SessionToLines()
	if err == nil {
		t.Log("Must be error")
		t.FailNow()
	}

	rs.Parameters.ColumnsSizes = []int{6, 12, 12, 8, 12, 4}
	_, err = rs.SessionToLines()
	if err == nil {
		t.Log("Must be error")
		t.FailNow()
	}

	rs.Parameters.FormatPattern = "%v%v%v%v%v%v"
	lines, err = rs.SessionToLines()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	var want int
	for _, c := range rs.Parameters.ColumnsSizes {
		want += c
	}
	for i := 4; i < len(lines); i++ {
		if len(lines[i]) != want {
			t.Errorf("Lenght is%v, must be %v", len(lines[i]), rs.Parameters.ColumnsSizes[i])
		}
	}
}

func TestSetDefaultParams(t *testing.T) {
	var rs *RaceSession
	var err error
	params := Parameters{InputfilePath: "/test/path", FontfilePath: "/test/path", OutputimagePath: "/test/path"}
	rs, err = NewRaceSession(params)
	if err != nil {
		t.Log(err)

		t.FailNow()
	}
	rs.setDefaultParams()
	p := rs.Parameters

	if p.FooterImagePath != defaults.FooterImagePath {
		t.Fail()
	}
	if p.FontSize != defaults.FontSize {
		t.Fail()
	}
	if p.FormatPattern != defaults.FormatPattern {
		t.Fail()
	}
	if !reflect.DeepEqual(p.ColumnsSizes, defaults.ColumnsSizes) {
		t.Fail()
	}
}

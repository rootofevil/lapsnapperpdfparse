package pdfParse

import (
	"os"
	"path"
	"reflect"
	"testing"
)

var contentdir string = "ContentTest"
var fontfilename string = "arial.ttf"
var footerfilename string = "logo.jpg"
var inputfilename string
var outfilename string

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
	// session, err := ReadPdf(inputfile)
	// lines := sessionToLines(session)
	// prepareImage(lines, postimage, fontfile, picture)

	inputfilename = "0433 Race  Need4Speed.pdf"
	outfilename = "out.jpg"
	outpath := path.Join(contentdir, outfilename)
	params := Parameters{InputfilePath: path.Join(contentdir, inputfilename),
		FontfilePath:    path.Join(contentdir, fontfilename),
		FooterImagePath: path.Join(contentdir, footerfilename),
		OutputimagePath: outpath}

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

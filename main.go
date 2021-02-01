package pdfParse

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/dcu/pdf"
	"github.com/golang/freetype"
	"github.com/hqbobo/text2pic"
)

type RaceSession struct {
	Type              string
	Name              string
	LapCount          string
	Started           string
	Ended             string
	TimeAttackResults []TimeAttack
	lines             []string
	text              string
	Parameters        Parameters
}

type TimeAttack struct {
	Position  int64
	Driver    string
	BestTime  string
	Dif       string
	TotalTime string
	Laps      int64
}

type Parameters struct {
	InputfilePath   string
	OutputimagePath string
	FontfilePath    string
	FooterImagePath string
	FontSize        float64
	FormatPattern   string
	ColumnsSizes    []int
}

var defaults Parameters = Parameters{
	FontSize:      6,
	FormatPattern: "%v%v%v%v%v%v",
	ColumnsSizes:  []int{6, 12, 12, 8, 12, 4},
}

func NewRaceSession(param Parameters) (*RaceSession, error) {
	if param.InputfilePath == "" {
		return nil, fmt.Errorf("Parameter %v is mandatory", "InputfilePath")
	}
	if param.FontfilePath == "" {
		return nil, fmt.Errorf("Parameter %v is mandatory", "FontfilePath")
	}
	if param.OutputimagePath == "" {
		return nil, fmt.Errorf("Parameter %v is mandatory", "OutputimagePath")
	}
	RaceSession := RaceSession{Parameters: param}
	return &RaceSession, nil
}

func (rs *RaceSession) setDefaultParams() {
	def := reflect.ValueOf(defaults)
	v := reflect.ValueOf(&rs.Parameters)
	elem := v.Elem()

	for i := 0; i < elem.NumField(); i++ {
		name := elem.Type().Field(i).Name
		field := elem.Field(i)

		if field.IsZero() && field.IsValid() {
			if toSet := def.FieldByName(name); toSet.IsValid() && field.CanSet() {
				field.Set(toSet)
			}

		}
	}
}

func (rs RaceSession) PdfToImage() (err error) {
	err = rs.ReadPdf()
	if err != nil {
		return
	}

	err = rs.GenerateImage()
	if err != nil {
		return
	}
	return nil
}

func (rs *RaceSession) ReadPdf() error {
	if len(rs.TimeAttackResults) != 0 {
		return fmt.Errorf("Session exist")
	}
	f, r, err := pdf.Open(rs.Parameters.InputfilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	// rs := RaceSession{}
	p := r.Page(2)
	rows, _ := p.GetTextByRow()
	if len(rows) >= 4 {
		rs.Type = rows[1].Content[0].S
		rs.Name = editString(rows[2].Content[0].S)    //strings.Trim(strings.Split(rows[2].Content[0].S, ":")[1], " ")
		rs.Started = editString(rows[2].Content[1].S) //strings.Trim(strings.Split(rows[2].Content[1].S, ":")[1], " ")
		rs.LapCount = editString(rows[3].Content[0].S)
		rs.Ended = editString(rows[3].Content[1].S)
	}
	for i := 6; i < len(rows); i++ {
		if r := rows[i]; len(r.Content) >= 6 {
			pos, err := strconv.ParseInt(strings.Trim(r.Content[0].S, "."), 10, 32)
			if err != nil {
				return err
			}
			driver := strings.Trim(r.Content[1].S, " ")
			best := strings.Trim(r.Content[2].S, " ")
			dif := strings.Trim(r.Content[3].S, " ")
			total := strings.Trim(r.Content[4].S, " ")
			laps, err := strconv.ParseInt(strings.Trim(r.Content[5].S, "."), 10, 32)
			if err != nil {
				return err
			}
			timeData := TimeAttack{Position: pos, Driver: driver, BestTime: best, Dif: dif, TotalTime: total, Laps: laps}
			rs.TimeAttackResults = append(rs.TimeAttackResults, timeData)
		}
	}
	return nil
}

func (rs *RaceSession) SessionToText() *RaceSession {
	text := fmt.Sprintf("%v\nStarted: %v\nEnded: %v\n\nPosition\tCar\t\tBest time\tDif\tTotal time\tLaps\n", rs.Type, rs.Started, rs.Ended)
	for _, r := range rs.TimeAttackResults {
		line := fmt.Sprintf("%v\t\t%v\t%v\t%v\t%v\t%v\n", r.Position, r.Driver, r.BestTime, r.Dif, r.TotalTime, r.Laps)
		text += line
	}
	rs.text = text
	return rs
}

func (rs *RaceSession) SessionToLines() *RaceSession {
	var lines []string
	columnsSize := rs.Parameters.ColumnsSizes
	pattern := rs.Parameters.FormatPattern
	lines = append(lines, fmt.Sprintf("%v", rs.Type))
	lines = append(lines, fmt.Sprintf("Started: %v", rs.Started))
	lines = append(lines, fmt.Sprintf("Ended:  %v", rs.Ended))
	lines = append(lines, " ")

	headoftable := fmt.Sprintf(pattern,
		incLen("Pos.", columnsSize[0]),
		incLen("Car", columnsSize[1]),
		incLen("Best time", columnsSize[2]),
		incLen("Dif", columnsSize[3]),
		incLen("Total time", columnsSize[4]),
		incLen("Laps", columnsSize[5]),
	)
	lines = append(lines, headoftable)
	for _, r := range rs.TimeAttackResults {

		line := fmt.Sprintf(pattern,
			incLen(fmt.Sprint(r.Position), columnsSize[0]),
			incLen(r.Driver, columnsSize[1]),
			incLen(r.BestTime, columnsSize[2]),
			incLen(r.Dif, columnsSize[3]),
			incLen(r.TotalTime, columnsSize[4]),
			incLen(fmt.Sprint(r.Laps), columnsSize[5]))
		lines = append(lines, line)
	}
	rs.lines = lines
	return rs
}

func (rs *RaceSession) GenerateImage() (err error) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(rs.Parameters.FontfilePath)
	if err != nil {
		return
	}
	//produce the fonttype
	f, err := freetype.ParseFont(fontBytes)

	if err != nil {
		return
	}
	text := rs.SessionToLines().lines
	pic := text2pic.NewTextPicture(text2pic.Configure{Width: 720, BgColor: text2pic.ColorBlack})
	var fontsize float64 = 8
	if rs.Parameters.FontSize != 0 {
		fontsize = rs.Parameters.FontSize
	}

	pic.AddTextLine(" ", fontsize, f, text2pic.ColorBlack, text2pic.Padding{})
	for _, l := range text {
		pic.AddTextLine(l, fontsize, f, text2pic.ColorWhite, text2pic.Padding{
			Left:      40,
			Right:     20,
			Bottom:    0,
			Top:       0,
			LineSpace: 0,
		})
	}
	pic.AddTextLine(" ", fontsize, f, text2pic.ColorBlack, text2pic.Padding{})

	if rs.Parameters.FooterImagePath != "" {
		file, err := os.Open(rs.Parameters.FooterImagePath)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
		pic.AddPictureLine(file, text2pic.Padding{Bottom: 40})
	}

	outFile, err := os.Create(rs.Parameters.OutputimagePath)
	if err != nil {
		return
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	//produce the output
	err = pic.Draw(b, text2pic.TypeJpeg)
	if err != nil {
		log.Print(err.Error())
	}
	e := b.Flush()
	if e != nil {
		fmt.Println(e)
	}
	return
}

func editString(s string) (res string) {
	raw := strings.Split(s, ":")

	for i := 1; i < len(raw); i++ {
		if i != 1 {
			res += ":"
		}
		res += strings.Trim(raw[i], " ")
	}
	return
}

func incLen(s string, l int) string {
	if len(s) > l {
		return s
	}

	for len(s) < l {
		s += " "
	}
	return s
}

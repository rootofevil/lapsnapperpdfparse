package pdfParse

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
}

type TimeAttack struct {
	Position  int64
	Driver    string
	BestTime  string
	Dif       string
	TotalTime string
	Laps      int64
}

func PdfToImage(inputfile, postimage, fontfile, picture string) (err error) {
	session, err := ReadPdf(inputfile)
	if err != nil {
		return
	}

	// message := sessionToText(session)
	// log.Print(message)
	lines := SessionToLines(session)
	PrepareImage(lines, postimage, fontfile, picture)
	return
}

func ReadPdf(path string) (RaceSession, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		log.Println()
	}
	defer f.Close()
	timeAttack := RaceSession{}
	p := r.Page(2)
	rows, _ := p.GetTextByRow()
	if len(rows) >= 4 {
		timeAttack.Type = rows[1].Content[0].S
		timeAttack.Name = editString(rows[2].Content[0].S)    //strings.Trim(strings.Split(rows[2].Content[0].S, ":")[1], " ")
		timeAttack.Started = editString(rows[2].Content[1].S) //strings.Trim(strings.Split(rows[2].Content[1].S, ":")[1], " ")
		timeAttack.LapCount = editString(rows[3].Content[0].S)
		timeAttack.Ended = editString(rows[3].Content[1].S)
	}
	for i := 6; i < len(rows); i++ {
		if r := rows[i]; len(r.Content) >= 6 {
			pos, err := strconv.ParseInt(strings.Trim(r.Content[0].S, "."), 10, 32)
			if err != nil {
				log.Println(err)
			}
			driver := strings.Trim(r.Content[1].S, " ")
			best := strings.Trim(r.Content[2].S, " ")
			dif := strings.Trim(r.Content[3].S, " ")
			total := strings.Trim(r.Content[4].S, " ")
			laps, err := strconv.ParseInt(strings.Trim(r.Content[5].S, "."), 10, 32)
			if err != nil {
				log.Println(err)
			}
			timeData := TimeAttack{Position: pos, Driver: driver, BestTime: best, Dif: dif, TotalTime: total, Laps: laps}
			timeAttack.TimeAttackResults = append(timeAttack.TimeAttackResults, timeData)
		}
	}
	return timeAttack, nil
}

func SessionToText(session RaceSession) string {
	text := fmt.Sprintf("%v\nStarted: %v\nEnded: %v\n\nPosition\tCar\t\tBest time\tDif\tTotal time\tLaps\n", session.Type, session.Started, session.Ended)
	for _, r := range session.TimeAttackResults {
		line := fmt.Sprintf("%v\t\t%v\t%v\t%v\t%v\t%v\n", r.Position, r.Driver, r.BestTime, r.Dif, r.TotalTime, r.Laps)
		text += line
	}
	return text
}

func SessionToLines(session RaceSession) []string {
	var lines []string
	lines = append(lines, fmt.Sprintf("%v", session.Type))
	lines = append(lines, fmt.Sprintf("Started: %v", session.Started))
	lines = append(lines, fmt.Sprintf("Ended:  %v", session.Ended))
	lines = append(lines, " ")
	lines = append(lines, "Pos. Car             Best time     Dif         Total time   Laps")
	for _, r := range session.TimeAttackResults {
		line := fmt.Sprintf("%v.     %v   %v   %v   %v   %v", r.Position, r.Driver, r.BestTime, r.Dif, r.TotalTime, r.Laps)
		lines = append(lines, line)
	}
	return lines
}

func PrepareImage(text []string, out string, fontpath string, imagepath string) {
	// Read the font data.
	fontBytes, err := ioutil.ReadFile(fontpath)
	if err != nil {
		log.Println(err)
		return
	}
	//produce the fonttype
	f, err := freetype.ParseFont(fontBytes)

	if err != nil {
		log.Println(err)
		return
	}
	pic := text2pic.NewTextPicture(text2pic.Configure{Width: 720, BgColor: text2pic.ColorBlack})
	pic.AddTextLine(" ", 8, f, text2pic.ColorBlack, text2pic.Padding{})
	for _, l := range text {
		pic.AddTextLine(l, 6, f, text2pic.ColorWhite, text2pic.Padding{
			Left:      40,
			Right:     20,
			Bottom:    0,
			Top:       0,
			LineSpace: 0,
		})
	}
	pic.AddTextLine(" ", 6, f, text2pic.ColorBlack, text2pic.Padding{})
	file, err := os.Open(imagepath)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	pic.AddPictureLine(file, text2pic.Padding{Bottom: 40})

	outFile, err := os.Create(out)
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

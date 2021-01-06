package pdfParse

import (
	"log"
	"strconv"
	"strings"

	"github.com/dcu/pdf"
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

// func main() {
// 	content, err := ReadPdf(os.Args[1]) // Read local pdf file
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("%+v", content)
// }

func ReadPdf(path string) (RaceSession, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		log.Fatal()
	}
	defer f.Close()
	// totalPage := r.NumPage()
	// fmt.Println("Total pages:", totalPage)
	// for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
	// 	p := r.Page(pageIndex)
	// 	if p.V.IsNull() {
	// 		continue
	// 	}

	// 	rows, _ := p.GetTextByRow()
	// 	for _, row := range rows {

	// 		println(">>>> row: ", row.Position)
	// 		for _, word := range row.Content {
	// 			fmt.Println(word.S)
	// 		}
	// 	}
	// }
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
			// fmt.Printf("%#v %#v %#v %#v %#v %#v\n", pos, driver, best, dif, total, laps)
		}
	}
	// for _, row := range rows {
	// 	fmt.Println(">>>> row:", row.Position)
	// 	for _, word := range row.Content {
	// 		fmt.Println(word.S)
	// 	}
	// }
	return timeAttack, nil
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

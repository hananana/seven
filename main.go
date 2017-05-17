package main

import (
	"encoding/csv"
	"fmt"
	. "gopkg.in/ahmetb/go-linq.v3"
	"gopkg.in/urfave/cli.v2"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type SevenData struct {
	Name         string
	TotalScore   float64
	Rank         float64
	Civilization string
	GameId       int
}

func main() {
	app := cli.NewApp()
	app.Action = func(c *cli.Context) error {
		lines, err := readLine(c.Args().First())
		if err != nil {
			return err
		}
		data, err := parse(lines)
		if err != nil {
			return err
		}

		distinctedByName := make([]SevenData, 0)
		From(data).DistinctByT(func(d SevenData) string {
			return d.Name
		}).ToSlice(&distinctedByName)

		longestName := From(distinctedByName).SelectT(func(d SevenData) string {
			return d.Name
		}).OrderByDescendingT(func(s string) int {
			return utf8.RuneCountInString(s)
		}).First().(string)
		longestNameCount := utf8.RuneCountInString(longestName)

		fmt.Println("----- プレイヤーの戦績 -----")
		for _, v := range distinctedByName {
			query := From(data).WhereT(func(d SevenData) bool {
				return d.Name == v.Name
			})
			averageRank := query.SelectT(func(d SevenData) float64 {
				return d.Rank
			}).Average()
			arStr := strconv.FormatFloat(averageRank, 'f', 1, 64)

			averageScore := query.SelectT(func(d SevenData) float64 {
				return d.TotalScore
			}).Average()
			asStr := strconv.FormatFloat(averageScore, 'f', 1, 64)

			diff := longestNameCount - utf8.RuneCountInString(v.Name)
			fmt.Println(v.Name + brank(diff) + ", " + asStr + ", " + arStr)
		}

		distinctedByCivil := make([]SevenData, 0)
		From(data).DistinctByT(func(d SevenData) string {
			return d.Civilization
		}).ToSlice(&distinctedByCivil)

		longestCivilName := From(distinctedByCivil).SelectT(func(d SevenData) string {
			return d.Civilization
		}).OrderByDescendingT(func(s string) int {
			return utf8.RuneCountInString(s)
		}).First().(string)

		longestCivilNameCount := utf8.RuneCountInString(longestCivilName)

		fmt.Println("----- 文明の戦績 -----")
		for _, v := range distinctedByCivil {
			query := From(data).WhereT(func(d SevenData) bool {
				return d.Civilization == v.Civilization
			})

			averageRank := query.SelectT(func(d SevenData) float64 {
				return d.Rank
			}).Average()
			arStr := strconv.FormatFloat(averageRank, 'f', 1, 64)

			averageScore := query.SelectT(func(d SevenData) float64 {
				return d.TotalScore
			}).Average()

			asStr := strconv.FormatFloat(averageScore, 'f', 1, 64)

			diff := longestCivilNameCount - utf8.RuneCountInString(v.Civilization)
			fmt.Println(v.Civilization + brank(diff) + ", " + asStr + ", " + arStr)
		}

		return nil
	}
	app.Run(os.Args)
}

func brank(c int) string {
	ret := ""
	for i := 0; i < c; i++ {
		ret += " "
	}
	return ret
}

func readLine(path string) ([][]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func parse(data [][]string) ([]SevenData, error) {
	ret := make([]SevenData, 0)
	for k, v := range data {
		if k != 0 {
			sd := SevenData{}
			sd.Name = strings.TrimSpace(v[1])
			sd.Civilization = strings.TrimSpace(v[2])
			sd.TotalScore, _ = strconv.ParseFloat(strings.TrimSpace(v[3]), 32)
			sd.Rank, _ = strconv.ParseFloat(strings.TrimSpace(v[4]), 32)
			ret = append(ret, sd)
		}
	}
	return ret, nil
}

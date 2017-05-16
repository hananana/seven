package main

import (
	"encoding/csv"
	"fmt"
	"gopkg.in/urfave/cli.v2"
	"os"
	"strconv"
	"strings"
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

		names, err := distinctName(data)
		if err != nil {
			return err
		}

		fmt.Println(data)

		return nil
	}
	app.Run(os.Args)
}

func distinctName(data []SevenData) ([]string, error) {
	ret := make([]string, 0)
	for _, v := range data {
		ret = append(ret, v)
	}
	return ret, nil
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

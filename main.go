package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

type loginData struct {
	isLogin bool
	logTime time.Time
}

func main() {
	const path = `Data/logins.csv`
	rawData, rawTimes := loadData(path)
	loginList := separateTimes(rawData, rawTimes)
	timePlayed := calculateTimes(loginList)
	writeToFile(timePlayed)
}

func loadData(path string) ([]string, []string) {
	fmt.Println("Loading data")
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	raw, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	const messageLine = 3
	const timeLine = 2
	times, data := []string{}, []string{}
	for _, line := range raw {
		times = append(times, line[timeLine])
		data = append(data, line[messageLine])
	}
	return data, times
}

func separateTimes(data, times []string) map[string][]loginData {
	const dateFormat = "1/2/2006 15:04"
	const loggedInLayout = " logged in"
	const nameSeparator = "**"
	isLogins, playerNames, rfcTimes := []bool{}, []string{}, []time.Time{}
	logins := make(map[string][]loginData)

	for i := 0; i < len(data); i++ {
		if len(data[i]) > len(loggedInLayout) {
			loginType := data[i][len(data[i])-len(loggedInLayout):] == loggedInLayout
			isLogins = append(isLogins, loginType)
		} else {
			continue
		}

		if strings.Contains(data[i], nameSeparator) {
			nameStart := strings.Index(data[i], nameSeparator) + len(nameSeparator)
			nameEnd := strings.LastIndex(data[i], nameSeparator)
			playerNames = append(playerNames, data[i][nameStart:nameEnd])
		} else {
			continue
		}

		formattedTime, err := time.Parse(dateFormat, times[i])
		if err != nil {
			continue
		}
		rfcTimes = append(rfcTimes, formattedTime)

		loginInstance := loginData{
			isLogin: isLogins[i],
			logTime: rfcTimes[i],
		}
		logins[playerNames[i]] = append(logins[playerNames[i]], loginInstance)
	}
	return logins
}

func calculateTimes(loginList map[string][]loginData) map[string]time.Duration {
	playerTimes := make(map[string]time.Duration)

	for player, logins := range loginList {
		for i := 1; i < len(logins); i++ {
			if !logins[i].isLogin && logins[i-1].isLogin {
				timeDifference := logins[i].logTime.Sub(logins[i-1].logTime)
				playerTimes[player] = playerTimes[player] + timeDifference
			}
		}
	}
	return playerTimes
}

func writeToFile(timePlayed map[string]time.Duration) {
	const fileName = "timeplayed.txt"
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for name, time := range timePlayed {
		_, err := fmt.Fprintln(file, name+" "+fmt.Sprint(time.Minutes()))
		if err != nil {
			panic(err)
		}
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully saved data to " + fileName)
}

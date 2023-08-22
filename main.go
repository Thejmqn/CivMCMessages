package main

import (
	"bufio"
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
	const loginPath = `Data/logins.csv`
	const messagePath = `Data/chat.csv`
	const names = "Data/deluvianames.txt"

	validNames := readNames(names)
	rawData, rawTimes := loadLoginData(loginPath, validNames)
	rawMessages := loadMessageData(messagePath, validNames)
	loginList := separateTimes(rawData, rawTimes)
	playerMessages := messageCount(rawMessages)
	timePlayed := calculateTimes(loginList)
	writeToFile(timePlayed, playerMessages)
}

func loadLoginData(path string, validNames []string) ([]string, []string) {
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
		if validNames == nil || stringsContain(validNames, line[messageLine]) {
			times = append(times, line[timeLine])
			data = append(data, line[messageLine])
		}
	}
	return data, times
}

func loadMessageData(path string, validNames []string) []string {
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
	messages := []string{}
	for _, line := range raw {
		if validNames == nil || stringsContain(validNames, line[messageLine]) {
			messages = append(messages, line[messageLine])
		}
	}
	return messages
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

func messageCount(messages []string) map[string]int {
	const startSeparator = "**["
	const endSeparator = "]**"
	playerMessageCounter := make(map[string]int)
	for _, message := range messages {
		if strings.Contains(message, startSeparator) && strings.Contains(message, endSeparator) {
			nameStart := strings.Index(message, startSeparator) + len(startSeparator)
			nameEnd := strings.LastIndex(message, endSeparator)
			name := message[nameStart:nameEnd]
			playerMessageCounter[name]++
		}
	}
	return playerMessageCounter
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

func readNames(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func writeToFile(timePlayed map[string]time.Duration, messages map[string]int) {
	const fileName = "timeplayed.txt"
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for name, time := range timePlayed {
		_, err := fmt.Fprintln(file, name+" "+fmt.Sprint(time.Minutes())+" "+fmt.Sprint(messages[name]))
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

func stringsContain(splice []string, key string) bool {
	for _, element := range splice {
		if strings.Contains(key, element) {
			return true
		}
	}
	return false
}

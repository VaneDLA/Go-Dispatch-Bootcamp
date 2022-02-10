package service

import (
	"Go-Dispatch-Bootcamp/types"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

type translatorService struct{}

func New() *translatorService {
	log.Println("In service | constructor")

	return &translatorService{}
}

func (ts *translatorService) readCsvFromFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("can not open file")
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("read file error")
	}
	return records, nil
}

func (ts *translatorService) GetUsers() (*[]types.User, error) {
	records, err := ts.readCsvFromFile("data.csv")
	if err != nil {
		return nil, err
	}

	var users []types.User

	for _, line := range records {
		id, err := strconv.Atoi(line[0])

		if err != nil {
			return nil, errors.New(fmt.Sprintf("Id '%v' is not a number", line[0]))
		}

		users = append(users, types.User{
			Id:         id,
			Username:   line[1],
			Identifier: line[2],
			FirstName:  line[3],
			LastName:   line[4],
		})
	}

	return &users, nil
}

func (ts *translatorService) GetUsersMap() (map[int]types.User, error) {
	records, err := ts.readCsvFromFile("data.csv")
	if err != nil {
		return nil, err
	}

	users := make(map[int]types.User, len(records))

	for _, line := range records {
		id, err := strconv.Atoi(line[0])

		if err != nil {
			return nil, errors.New(fmt.Sprintf("Id '%v' is not a number", line[0]))
		}

		users[id] = types.User{
			Id:         id,
			Username:   line[1],
			Identifier: line[2],
			FirstName:  line[3],
			LastName:   line[4],
		}
	}

	return users, nil
}

func (ts *translatorService) GetFeedUsers() (*[]types.FeedUser, error) {
	records, err := ts.readCsvFromFile("feed.csv")
	if err != nil {
		return nil, err
	}

	var users []types.FeedUser

	for i, line := range records {
		if i == 0 {
			continue
		}

		users = append(users, types.FeedUser{
			Username:   line[0],
			Identifier: line[1],
			FirstName:  line[2],
			LastName:   line[3],
		})
	}

	return &users, nil
}
package repository

import (
	"bufio"
	"fmt"
	"lab3/models"
	"os"
	"sort"
	"strconv"
	"strings"
)

// FileRepo define repo
type FileRepo struct {
	reader *os.File
	writer *os.File
}

type ObjectSlice []*models.Сrossroad

func (s ObjectSlice) Len() int {
	return len(s)
}

func (s ObjectSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ObjectSlice) Less(i, j int) bool {
	return s[i].CrossroadID < s[j].CrossroadID
}

// NewFileRepo used to init FileRepo
func NewFileRepo(reader *os.File, writer *os.File) *FileRepo {
	return &FileRepo{reader: reader, writer: writer}
}

func (fr *FileRepo) GetCrossroads() ([]*models.Сrossroad, *models.Condition, error) {
	nodes := make(map[int]*models.Сrossroad)
	var condition models.Condition
	fileScanner := bufio.NewScanner(fr.reader)

	// read first row
	if fileScanner.Scan() {
		firstLine := fileScanner.Text()
		strArr := strings.Split(firstLine, " ")
		n, err := strconv.Atoi(strArr[0])
		if err != nil {
			return nil, nil, fmt.Errorf("error while read N - %v", err)
		}
		m, err := strconv.Atoi(strArr[1])
		if err != nil {
			return nil, nil, fmt.Errorf("error while read M - %v", err)
		}
		condition.N = n
		condition.M = m
	}
	var rows []string
	for fileScanner.Scan() {
		rows = append(rows, fileScanner.Text())
	}
	for i, row := range rows {
		if i == len(rows)-1 {
			break
		}
		// first vertex
		var firstCrossroad models.Сrossroad
		// second vertex
		var secondCrossroad models.Сrossroad

		strArr := strings.Split(row, " ")
		firstCrossroadID, err := strconv.Atoi(strArr[0])
		if err != nil {
			return nil, nil, fmt.Errorf("error while read СrossroadID - %v", err)
		}
		secondCrossroadID, err := strconv.Atoi(strArr[1])
		if err != nil {
			return nil, nil, fmt.Errorf("error while read СrossroadQ - %v", err)
		}
		roadTime, err := strconv.Atoi(strArr[2])
		if err != nil {
			return nil, nil, fmt.Errorf("error while read СrossroadRoad - %v", err)
		}
		firstCrossroad.CrossroadID = firstCrossroadID - 1
		firstCrossroad.Dist = 1e9
		secondCrossroad.CrossroadID = secondCrossroadID - 1
		secondCrossroad.Dist = 1e9
		// first edge
		firstRoad := models.Road{StartCrossroad: &firstCrossroad, EndCrossroad: &secondCrossroad, RoadTime: roadTime}
		// second edge
		secondRoad := models.Road{StartCrossroad: &secondCrossroad, EndCrossroad: &firstCrossroad, RoadTime: roadTime}

		// add first edge to first vertex
		firstCrossroad.CrossroadRoads = append(firstCrossroad.CrossroadRoads, &firstRoad)

		// add second edge to second vertex
		secondCrossroad.CrossroadRoads = append(secondCrossroad.CrossroadRoads, &secondRoad)

		// check is first vertex exists
		// if not exists add vertex to map
		// if exists add edge to current vertex
		existFirstNode, ok := nodes[firstCrossroad.CrossroadID]

		if !ok {
			nodes[firstCrossroad.CrossroadID] = &firstCrossroad
		} else {
			existFirstNode.CrossroadRoads = append(existFirstNode.CrossroadRoads, &firstRoad)
		}

		// check is second vertex exists
		// if not exists add vertex to map
		// if exists add edge to current vertex
		existSecondNode, ok := nodes[secondCrossroad.CrossroadID]
		if !ok {
			nodes[secondCrossroad.CrossroadID] = &secondCrossroad
		} else {
			existSecondNode.CrossroadRoads = append(existSecondNode.CrossroadRoads, &secondRoad)
		}
	}
	strArr := strings.Split(rows[len(rows)-1], " ")
	s, err := strconv.Atoi(strArr[0])
	if err != nil {
		return nil, nil, fmt.Errorf("error while read S - %v", err)
	}
	f, err := strconv.Atoi(strArr[1])
	if err != nil {
		return nil, nil, fmt.Errorf("error while read F - %v", err)
	}
	q, err := strconv.Atoi(strArr[2])
	if err != nil {
		return nil, nil, fmt.Errorf("error while read Q - %v", err)
	}
	condition.S = s
	condition.F = f
	condition.Q = q

	// add time of crossing current crossroad to road time
	var nodesArr []*models.Сrossroad
	for _, node := range nodes {
		for _, road := range node.CrossroadRoads {
			road.RoadTime = road.RoadTime + condition.Q*len(node.CrossroadRoads)
		}
		nodesArr = append(nodesArr, node)
	}
	sort.Sort(ObjectSlice(nodesArr))
	return nodesArr, &condition, nil
}

func (fr *FileRepo) PrintResult(result string) error {
	file, err := os.Create("output.txt")
	if err != nil {
		return fmt.Errorf("error creating file - %v", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(result)
	if err != nil {
		return fmt.Errorf("error writing to file - %v", err)
	}

	writer.Flush()
	return nil
}

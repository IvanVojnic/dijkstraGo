package models

type Сrossroad struct {
	CrossroadID    int
	CrossroadRoads []*Road
	Time           int
	Dist           int
}

package models

type BullpenRoleProb struct {
	Inning           int
	RunsScoredGame   int
	RunsScoredInning int
	Role             int
	PullProbability  float64
}

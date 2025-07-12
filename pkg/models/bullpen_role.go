package models

// Define this struct somewhere, maybe in utils.go:
type BullpenRoleProb struct {
	Inning    int
	RunDiff   int
	RunnersOn int
	Role      int
	Prob      float64
}

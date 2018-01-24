package model

// FindWeightPlates calculates the number of each weight plate needed for a given weight.
func FindWeightPlates(weight int, barWeight int) WeightPlate {

	var answer = WeightPlate{Plate45: 0, Plate35: 0, Plate25: 0, Plate10: 0, Plate05: 0, Plate02: 0}
	remaining := weight - barWeight
	perSide := remaining / 2

	for perSide >= 45 {
		answer.Plate45++
		perSide -= 45
	}
	for perSide >= 35 {
		answer.Plate35++
		perSide -= 35
	}
	for perSide >= 25 {
		answer.Plate25++
		perSide -= 25
	}
	for perSide >= 10 {
		answer.Plate10++
		perSide -= 10
	}
	for perSide >= 5 {
		answer.Plate05++
		perSide -= 5
	}
	if perSide > 2 {
		answer.Plate02++
	}

	return answer
}

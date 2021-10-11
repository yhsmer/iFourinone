package main

import (
	"iFourinone/class"
)

func main() {
	var contractor class.Contractor
	m := map[string]string{"filepath": "./in"}
	contractor.ScheduleWork(m)

	// 工头之间的串行处理
	//contractor.ToNext(class.Contractor{Name: "TwoCtor"},m).ToNext(class.Contractor{Name: "ThreeCtor"},m)

}

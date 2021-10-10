package main

import (
	"iFourinone/class"
)

func main() {
	var contractor class.Contractor
	m := map[string]string{"filepath": "./in"}
	contractor.ScheduleWork(m)
}

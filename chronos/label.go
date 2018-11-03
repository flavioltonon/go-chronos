package chronos

type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LabelSpec struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func SetColorToLabel(name string) string {
	if name == "Overdue" {
		return "4b21c6"
	}
	// number, err := strconv.Atoi(strings.Split(name, " ")[1])
	// if err != nil {
	// 	return "ffffff"
	// }

	// switch strings.Split(name, " ")[2] {
	// case "horas":
	// 	return "990002"
	// case "dias":
	// 	if number > 15 {
	// 		return "0e8a16"
	// 	}
	// 	if number > 7 {
	// 		return "fffa50"
	// 	}
	// 	if number > 3 {
	// 		return "ffa550"
	// 	}
	// 	if number > 1 {
	// 		return "c60003"
	// 	}
	// 	return "ffffff"
	// default:
	// 	return "ffffff"
	// }
	return "ffffff"
}

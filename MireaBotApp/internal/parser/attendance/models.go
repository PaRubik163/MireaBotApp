package attendance

import "fmt"

type Subject struct {
	Name   string
	Points float64
}

func (s *Subject) FromMapToStruct(res []map[string]interface{}) []Subject {
	var subjects []Subject
	for _, m := range res {
		name, ok1 := m["name"].(string)
		points, ok2 := m["points"].(float64) // или int, если баллы целые

		if !ok1 || !ok2 {
			fmt.Println("Некорректные данные в мапе:", m)
			continue
		}

		subjects = append(subjects, Subject{
			Name:   name,
			Points: points,
		})
	}

	return subjects
}

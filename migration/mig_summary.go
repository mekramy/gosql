package migration

type Migrated struct {
	Name  string `db:"name"`
	Stage string `db:"stage"`
}

type Summary []Migrated

func (s Summary) IsEmpty() bool {
	return len(s) == 0
}

func (s Summary) Names() []string {
	result := make([]string, 0)
	for _, migration := range s {
		result = append(result, migration.Name)
	}
	return result
}

func (s Summary) GroupByStage() map[string][]string {
	result := make(map[string][]string)
	for _, file := range s {
		result[file.Stage] = append(result[file.Stage], file.Name)
	}
	return result
}

func (s Summary) GroupByFile() map[string][]string {
	result := make(map[string][]string)
	for _, file := range s {
		result[file.Name] = append(result[file.Name], file.Stage)
	}
	return result
}

func (s Summary) ForStage(stage string) Summary {
	result := make(Summary, 0)
	for _, file := range s {
		if file.Stage == stage {
			result = append(result, Migrated{
				Name:  file.Name,
				Stage: file.Stage,
			})
		}
	}
	return result
}

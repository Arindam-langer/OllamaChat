package main

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}


// Initial State of Application

funct initialModel() model{
	return model{
		choices []string{"carrot","banana","rabbit meat"}
		selected map[int]struct{}
	}
}

func (m model) Init() tea.Cmd {
	return nil
} 

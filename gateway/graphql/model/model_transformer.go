package model

func (c *Categories) IDs() []int {
	ids := make([]int, 0, len(c.Data))
	for _, e := range c.Data {
		ids = append(ids, e.ID)
	}
	return ids
}

func (c *Tags) IDs() []int {
	ids := make([]int, 0, len(c.Data))
	for _, e := range c.Data {
		ids = append(ids, e.ID)
	}
	return ids
}

func (c *Games) IDs() []int {
	ids := make([]int, 0, len(c.Data))
	for _, e := range c.Data {
		ids = append(ids, e.ID)
	}
	return ids
}

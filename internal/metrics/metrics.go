package metrics

type MemStorage struct {
	gauge   map[string]float32
	counter map[string]int
}

func NewMemStorage() MemStorage {
	return MemStorage{}
}

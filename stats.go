package cgminer

import "encoding/json"

// S9 - get struct fields only related to S9
func (s *GenericStats) S9() (*StatsS9, error) {
	raw, _ := json.Marshal(s)
	result := &StatsS9{}
	err := json.Unmarshal(raw, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// S7 - get struct fields only related to S7
func (s *GenericStats) S7() (*StatsS7, error) {
	raw, _ := json.Marshal(s)
	result := &StatsS7{}
	err := json.Unmarshal(raw, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// L3 - get struct fields only related to L3+
func (s *GenericStats) L3() (*StatsL3, error) {
	raw, _ := json.Marshal(s)
	result := &StatsL3{}
	err := json.Unmarshal(raw, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// D3 - get struct fields only related to D3
func (s *GenericStats) D3() (*StatsD3, error) {
	raw, _ := json.Marshal(s)
	result := &StatsD3{}
	err := json.Unmarshal(raw, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// T9 - get struct fields only related to T9+
func (s *GenericStats) T9() (*StatsT9, error) {
	raw, _ := json.Marshal(s)
	result := &StatsT9{}
	err := json.Unmarshal(raw, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

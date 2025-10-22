package data_structure

type SimpleSet struct {
	key  string
	dict map[string]struct{}
}

func NewSimpleSet(key string) *SimpleSet {
	return &SimpleSet{
		key:  key,
		dict: make(map[string]struct{}),
	}
}

func (s *SimpleSet) Add(members ...string) int {
	added := 0
	for _, member := range members {
		if _, exist := s.dict[member]; !exist {
			s.dict[member] = struct{}{}
			added++
		}
	}
	return added
}

func (s *SimpleSet) Rem(members ...string) int {
	removed := 0
	for _, member := range members {
		if _, exist := s.dict[member]; exist {
			delete(s.dict, member)
			removed++
		}
	}
	return removed
}

func (s *SimpleSet) IsMember(member string) int {
	_, exist := s.dict[member]
	if exist {
		return 1
	}
	return 0
}

func (s *SimpleSet) Members() []string {
	member := make([]string, 0, len(s.dict))
	for k, _ := range s.dict {
		member = append(member, k)
	}
	return member
}

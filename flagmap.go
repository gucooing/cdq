package cdq

type FlagMap map[string]*FlagMapItem

type FlagMapItem struct {
	Value interface{}
}

func (f FlagMap) String(long string) string {
	i := f[long]
	if i == nil {
		return ""
	}
	s, ok := i.Value.(string)
	if !ok {
		return ""
	}
	return s
}

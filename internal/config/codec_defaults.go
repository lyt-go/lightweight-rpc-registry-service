package config

type CodecRule interface{ Allow(string) bool }

type PrefixRule struct{ Prefix string }

func (r *PrefixRule) Allow(name string) bool {
	if r == nil {
		return true
	}
	return len(name) >= len(r.Prefix) && name[:len(r.Prefix)] == r.Prefix
}

type CodecDefaults struct {
	Aliases   map[string]string
	Validator CodecRule
}

func LoadCodecDefaults() CodecDefaults {
	var rule *PrefixRule
	return CodecDefaults{Validator: rule}
}

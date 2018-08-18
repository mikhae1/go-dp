package config

// Env is a yaml config fromat
// Remote for remote commands
// Local for local commands
// General for both local and remote commands and merged to Local and Remote
type Env struct {
	Parent   string            `yaml:"parent"`
	Hidden   bool              `yaml:"hidden"`
	Template map[string]string `yaml:"template"`
	General  envTarget         `yaml:"general"`
	Remote   envTarget         `yaml:"remote"`
	Local    envTarget         `yaml:"local"`
	Targets  map[string]struct {
		General envTarget `yaml:"general"`
		Remote  envTarget `yaml:"remote"`
		Local   envTarget `yaml:"local"`
	} `yaml:"targets"`
}

// custom types for shorten yamls
type mapOrString map[string]string
type sliceOrString []string

type envTarget struct {
	Servers sliceOrString `yaml:"servers"`
	User    string        `yaml:"user"`
	Log     mapOrString   `yaml:"log"`
	Cmd     mapOrString   `yaml:"cmd"`
	Cat     mapOrString   `yaml:"cat"`
	Branch  string        `yaml:"branch"`
	URL     string        `yaml:"url"`
	Path    string        `yaml:"path"`
}

// implements the yaml.Unmarshaler interface for types that could be string or maps
func (m *mapOrString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection

	var mapValue map[string]string
	err := unmarshal(&mapValue)
	if err != nil {
		var stringValue string
		err := unmarshal(&stringValue)
		if err != nil {
			return err
		}

		*m = map[string]string{"0": stringValue}
	} else {
		*m = mapValue
	}
	return nil
}

// implements the yaml.Unmarshaler interface for types that could be string or array
func (s *sliceOrString) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	var sliceValue []string

	err = unmarshal(&sliceValue)
	if err == nil {
		*s = sliceValue
		return
	}

	var stringValue string
	if err = unmarshal(&stringValue); err != nil {
		return
	}

	*s = []string{stringValue}
	return
}

// // implements the yaml.Unmarshaler interface for flexible log types
// func (m *MapOrString) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
// 	// again, we have to hide it using a type indirection.
// 	type plain MapOrString

// 	if err := unmarshal((*plain)(er)); err != nil {
// 		return err
// 	}

// 	for _, key := range keyList {
// 		switch reflect.TypeOf(er.LogRaw).Kind() {
// 		case reflect.String:
// 			log.Println("string!")
// 			er.LogStr = er.LogRaw.(string)
// 		case reflect.Map:
// 			rawMap := reflect.ValueOf(er.LogRaw)
// 			er.LogMap = make(map[string]string)
// 			for _, v := range rawMap.MapKeys() {
// 				er.LogMap[v.Interface().(string)] = rawMap.MapIndex(v).Interface().(string)
// 			}
// 		}
// 	}

// 	return nil
// }

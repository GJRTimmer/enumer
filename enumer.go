package main

import "fmt"

// Arguments to format are:
//	[1]: type name
const stringValueToNameMap = `func %[1]sString(s string) (%[1]s, error) {
	if val, ok := _%[1]sNameToValue_map[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%%s does not belong to %[1]s values", s)
}
`

const listMethod = `
func %[1]sList() []%[1]s {
	list := make([]%[1]s, len(_%[1]sNameToValue_map))
	idx := 0
	for _, v := range _%[1]sNameToValue_map {
		list[idx] = v
		idx++
	}
	return list
}

func %[1]sListString() []string {
	list := make([]string, len(_%[1]sNameToValue_map))
	idx := 0
	for k := range _%[1]sNameToValue_map {
		list[idx] = k
		idx++
	}
	return list
}
`

const isValidmethod = `
func %[1]sIsValid(t %[1]s) bool {
	for _, v := range %[1]sList() {
		if t == v {
			return true
		}
	}
	return false
}
`

func (g *Generator) buildValueToNameMap(runs [][]Value, typeName string, runsThreshold int) {
	// At this moment, either "g.declareIndexAndNameVars()" or "g.declareNameVars()" has been called
	g.Printf("\nvar _%sNameToValue_map = map[string]%s{\n", typeName, typeName)
	thereAreRuns := len(runs) > 1 && len(runs) <= runsThreshold
	var n int
	var runID string
	for i, values := range runs {
		if thereAreRuns {
			runID = "_" + fmt.Sprintf("%d", i)
			n = 0
		} else {
			runID = ""
		}

		for _, value := range values {
			g.Printf("\t_%s_name%s[%d:%d]: %s,\n", typeName, runID, n, n+len(value.name), &value)
			n += len(value.name)
		}
	}
	g.Printf("}\n\n")
	g.Printf(stringValueToNameMap, typeName)
	g.Printf(listMethod, typeName)
	g.Printf(isValidmethod, typeName)
}

// Arguments to format are:
//	[1]: type name
const jsonMethods = `
func (i %[1]s) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *%[1]s) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("%[1]s should be a string, got %%s", data)
	}

	var err error
	*i, err = %[1]sString(s)
	return err
}
`

func (g *Generator) buildJSONMethods(runs [][]Value, typeName string, runsThreshold int) {
	g.Printf(jsonMethods, typeName)
}

// Arguments to format are:
//	[1]: type name
const yamlMethods = `
func (i %[1]s) MarshalYAML() (interface{}, error) {
	return i.String(), nil
}

func (i *%[1]s) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var err error
	*i, err = %[1]sString(s)
	return err
}
`

func (g *Generator) buildYAMLMethods(runs [][]Value, typeName string, runsThreshold int) {
	g.Printf(yamlMethods, typeName)
}

package form

import (
	"html/template"
	"reflect"
	"strings"
)

func fields(v interface{}, parents ...string) []field {
	rv := reflect.ValueOf(v)

	// Try to get the element, not an interface or pointer
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		panic("form: invalid value - only structs are supported")
	}

	t := rv.Type()
	ret := make([]field, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		var tmp reflect.Value
		if t.Field(i).Type.Kind() == reflect.Ptr {
			// if pointer, get a value not the pointer
			// to it
			tmp = reflect.New(t.Field(i).Type.Elem()).Elem()
		} else {
			tmp = reflect.New(t.Field(i).Type).Elem()
		}

		// If this is a struct, add its nested fields
		// and move on to the next field
		if tmp.Kind() == reflect.Struct {
			ret = append(ret, fields(rv.Field(i).Interface(), append(parents, t.Field(i).Name)...)...)
			continue
		}

		tags := parseTags(t.Field(i).Tag.Get("form"))
		nameStrs := append(parents, t.Field(i).Name)
		f := field{
			Name:        strings.Join(nameStrs, "."),
			Label:       t.Field(i).Name,
			Placeholder: t.Field(i).Name,
			Type:        "text",
			Value:       rv.Field(i).Interface(),
		}
		if v, ok := tags["name"]; ok {
			f.Name = v
		}
		if v, ok := tags["label"]; ok {
			f.Label = v
			// Will be overwritten if there is a placeholder set
			f.Placeholder = v
		}
		if v, ok := tags["placeholder"]; ok {
			f.Placeholder = v
		}
		if v, ok := tags["type"]; ok {
			f.Type = v
		}
		if v, ok := tags["id"]; ok {
			f.ID = v
		}
		if v, ok := tags["footer"]; ok {
			f.Footer = template.HTML(v)
		}
		if _, ok := tags["ignore"]; !ok {
			ret = append(ret, f)
		}
	}
	return ret
}

func parseTags(tags string) map[string]string {
	tags = strings.TrimSpace(tags)
	if len(tags) == 0 {
		return map[string]string{}
	}
	split := strings.Split(tags, ";")
	ret := make(map[string]string, len(split))
	for _, tag := range split {
		kv := strings.Split(tag, "=")
		if len(kv) < 2 {
			if kv[0] == "-" {
				return map[string]string{
					"ignore": "yes please",
				}
			}
			continue
		}
		k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		ret[k] = v
	}
	return ret
}

type field struct {
	Name        string
	Label       string
	Placeholder string
	Type        string
	ID          string
	Value       interface{}
	Footer      template.HTML
}

package form

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

func fields(v interface{}, name ...string) []field {
	rv := reflect.ValueOf(v)
	// If a nil pointer is passed in but has a type we can recover, but I
	// really should just panic and tell people to fix their shitty code.
	if rv.Type().Kind() == reflect.Ptr && rv.IsNil() {
		rv = reflect.New(rv.Type().Elem()).Elem()
	}
	// If we have a pointer or interface let's try to get the underlying
	// element
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		fmt.Println(rv.Kind())
		panic("invalid value; only structs are supported")
	}

	t := rv.Type()
	ret := make([]field, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		rf := rv.Field(i)
		// If this is a nil pointer, create a new instance of the element
		// type that it points to so we can recur easier.
		if t.Field(i).Type.Kind() == reflect.Ptr && rf.IsNil() {
			rf = reflect.New(t.Field(i).Type.Elem()).Elem()
		}

		// If this is a struct it has nested fields we need to add. The
		// simplest way to do this is to recursively call `fields` but
		// to provide the name of this struct field to be added as a prefix
		// to the fields.
		if rf.Kind() == reflect.Struct {
			ret = append(ret, fields(rf.Interface(), append(name, t.Field(i).Name)...)...)
			continue
		}

		// If we are still in this loop then we aren't dealing with a nested
		// struct and need to add the field. First we check to see if the
		// ignore tag is present, then we set default values, then finally
		// we overwrite defaults with any provided tags.
		tags := parseTags(t.Field(i).Tag.Get("form"))
		if _, ok := tags["-"]; ok {
			continue
		}

		fieldName := append(name, t.Field(i).Name)
		f := field{
			Name:        strings.Join(fieldName, "."),
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
			// DO NOT move this label check after the placeholder check or
			// this will cause issues.
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
			// Probably shouldn't be HTML but whatever.
			f.Footer = template.HTML(v)
		}

		// The "-" tag is special and signified a field to ignore
		// and not add to your list of fields.
		if _, ok := tags["-"]; !ok {
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
					"-": "this field is ignored",
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

package form

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"
	"regexp"
)


var selects map[string]map[string]interface{}
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z][0-9]+)")
var matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
var matchNumbers = regexp.MustCompile("([A-Za-z]+)([0-9]+)")

// valueOf is basically just reflect.ValueOf, but if the Kind() of the
// value is a pointer or interface it will try to get the reflect.Value
// of the underlying element, and if the pointer is nil it will
// create a new instance of the type and return the reflect.Value of it.
//
// This is used to make the rest of the fields function simpler.
func valueOf(v interface{}) reflect.Value {

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
	return rv
}

func RegisterSelect(key string,opts map[string]interface{}){

if(selects==nil){
	selects = make(map[string]map[string]interface{})
}

selects[key]=opts


}

func fields(v interface{}, names ...string) []field {
	rv := valueOf(v)
	if rv.Kind() != reflect.Struct {
		// We can't really do much with a non-struct type. I suppose this
		// could eventually support maps as well, but for now it does not.
		panic("invalid value; only structs are supported")
	}

	t := rv.Type()
	ret := make([]field, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		rf := rv.Field(i)
		// If this is a nil pointer, create a new instance of the element.
		// This could probably be done in a simpler way given that we
		// typically recur with this value, but this works so I'm letting it
		// be.
		if t.Field(i).Type.Kind() == reflect.Ptr && rf.IsNil() {
			rf = reflect.New(t.Field(i).Type.Elem()).Elem()
		}

		// If this is a struct it has nested fields we need to add. The
		// simplest way to do this is to recursively call `fields` but
		// to provide the name of this struct field to be added as a prefix
		// to the fields.
		
		if(!rf.CanInterface()){
		           continue
		}
		
		if rf.Kind() == reflect.Struct {
			
		
			
			ret = append(ret, fields(rf.Interface(), append(names, t.Field(i).Name)...)...)
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
		name := append(names, t.Field(i).Name)
		f := field{
			Name:        strings.Join(name, "."),
			Label:       fromCamelCase(t.Field(i).Name),
			Placeholder: fromCamelCase(t.Field(i).Name),
			Type:        "text",
			Value:       rv.Field(i).Interface(),
		}
		applyTags(&f, tags)
		ret = append(ret, f)
	}
	return ret
}


func fromCamelCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
	snake  = matchAllCap.ReplaceAllString(snake, "${1} ${2}")
	snake = matchNumbers.ReplaceAllString(snake, "${1} ${2}")
	return string(snake)
}


func applyTags(f *field, tags map[string]string) {
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

	var lns []string

	for k,v := range tags{

		lns = append(lns,fmt.Sprintf(`%s=%q`,k,v))

	}

	f.Attrs = template.HTMLAttr(strings.Join(lns," "))


	if(f.Type=="select" || f.Type=="checkbox"){

         if it,oks := selects[f.Name]; oks{

         	f.Items = it
            //this block allows us to set the select value as an output ie CA=California, f.Value is CA and f.SelectValue is California
         	for v,k := range it {
         		if(k==f.Value){
         			f.SelectValue = v
				}

			}


		}



		if v, ok := tags["select"]; ok {
			f.SelectType = template.HTMLAttr(v)
		}

	}

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
	SelectValue interface{}
	Footer      template.HTML
	Items       map[string]interface{}
	SelectType template.HTMLAttr
	Attrs template.HTMLAttr
}

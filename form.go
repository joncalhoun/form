package form

import (
	"github.com/Masterminds/sprig"
	"html/template"
	"io/ioutil"
	"strings"
	"time"
)

var basetmpl string

type Form struct {
	Tpl *template.Template
    selectMap map[string]map[string]interface{}
	Action string
	Method string
	Skip []string
}

func init() {

	basetmpl = `<div class="form-group">
    {{if ne .Type "hidden"}}
    <label class="form-label" {{with .ID}}for="{{.}}"{{end}}>
        {{.Label}}
    </label>
    {{end}}
    {{if eq .Type "textarea"}}
    <textarea {{.Attrs}} class="form-control" {{with .ID}}id="{{.}}"{{end}} name="{{.Name}}" rows="3" placeholder="{{.Placeholder}}">{{with .Value}}{{.}}{{end}}</textarea>
    {{else if eq .Type "checkbox" }}
    <input {{.Attrs}} type="{{.Type}}"  class="form-check-input" {{with .ID}}id="{{.}}"{{end}} name="{{.Name}}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}}>
    {{else if eq .Type "select" }}
    <select {{.Attrs}} class="form-control" {{with .ID}}id="{{.}}"{{end}} name="{{.Name}}" {{.SelectType}}>
        {{ $myval := .Value }}
        {{ if gt (len .Placeholder) 0 }}
        <option value="" >{{ .Placeholder }}</option>
        {{ end }}
        {{ range $v,$k := .Items}}
          <option {{ if eq $myval $k  }}selected="selected"{{end}}value="{{$k}}">{{$v}}</option>
        {{end}}
    </select>
    {{ else }}
    <input {{.Attrs}} type="{{.Type}}" class="form-control" {{with .ID}}id="{{.}}"{{end}} name="{{.Name}}" placeholder="{{.Placeholder}}" {{with .Value}}value="{{.}}"{{end}}>
    {{end}}
    {{with .Footer}}
    <small class="form-text text-muted"> {{.}} </small>
    {{end}}
</div>`

}

func New(pth ...string) (*Form,error){
    var frmstr string
    var p string

    if(len(pth)>0){

    	p = pth[0]
	}else{
		p = ""
	}
	frm, errf := ioutil.ReadFile(p)
	if errf != nil {
		frmstr = basetmpl
	}else{

		frmstr = string(frm)
	}

	tpl := template.Must(template.New("form").Funcs(sprig.FuncMap()).Funcs(template.FuncMap{
		"datetimelocal": func(val interface{}) string{
			return val.(time.Time).Format("2006-01-02T15:04")
		},
		"datetime": func(val interface{}) string{
			return val.(time.Time).Format("01/02/2006 15:04")
		},
	}).Parse(frmstr))


	return &Form{Tpl: tpl},errf

}

func (f *Form) SkipField(skip string){

  f.Skip = append(f.Skip,skip)
}

func (f *Form) Select(nm string,mp map[string]interface{}){

	if(f.selectMap==nil){
		f.selectMap = make(map[string]map[string]interface{})
	}

	f.selectMap[nm]=mp

}



func (f *Form) Render(v interface{}, errs ...error) (template.HTML, error) {

	fields := fields(v)

	errors := fieldErrors(errs)
	var html template.HTML
	for _, field := range fields {

		dump := false



		for _,sv := range f.Skip {


			if(sv == field.Name){
				dump = true
				break
			}
		}

		if(dump==true){continue}


		if(field.Type=="select" || field.Type=="checkbox" ){

			if it,oks := f.selectMap[field.Name]; oks{

				field.Items = it

				//this block allows us to set the select value as an output ie CA=California, f.Value is CA and f.SelectValue is California
				for v,k := range it {
					if(k==field.Value){
						field.SelectValue = v
					}

				}

			}

		}


	var sb strings.Builder
	f.Tpl.Funcs(template.FuncMap{
		"errors": func() []string {
		if errs, ok := errors[field.Name]; ok {
		return errs
		}
		return nil
		},
	})
	err := f.Tpl.Execute(&sb, field)
	if err != nil {
	return "", err
	}
	html = html + template.HTML(sb.String())
	}
	return html, nil

}
package builtin

import (
	"fmt"

	"github.com/trevex/kale/pkg/module"
	"github.com/trevex/kale/pkg/util"
	"go.starlark.net/starlark"
)

var (
	defaultString = ""
	defaultBool   = false
	defaultInt    = 0
	defaultFloat  = 0.0
)

func SchemaModule() starlark.Value {
	mod := &module.Module{}
	mod.SetKeyFunc(starlark.String("filename"), primitiveSchemaObjectFunc("filename", defaultString))
	mod.SetKeyFunc(starlark.String("string"), primitiveSchemaObjectFunc("string", defaultString))
	mod.SetKeyFunc(starlark.String("bool"), primitiveSchemaObjectFunc("bool", defaultBool))
	mod.SetKeyFunc(starlark.String("int"), primitiveSchemaObjectFunc("int", defaultInt))
	mod.SetKeyFunc(starlark.String("float"), primitiveSchemaObjectFunc("float", defaultFloat))
	return mod
}

type SchemaObject struct {
	Type     string
	Required bool
	Default  interface{}
}

func NewSchemaObject(name string, defaultValue interface{}) *SchemaObject {
	return &SchemaObject{
		Type:     name,
		Required: false,
		Default:  defaultValue,
	}
}

func SchemaObjectFromDict(dict *starlark.Dict) (*SchemaObject, error) {
	obj := &SchemaObject{}
	// type
	v, ok, err := dict.Get(starlark.String("type"))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("Field 'type' missing! Not a schema-object?")
	}
	obj.Type, ok = starlark.AsString(v)
	if !ok {
		return nil, fmt.Errorf("Field 'type' not a string! Not a schema-object?")
	}
	// required
	v, ok, err = dict.Get(starlark.String("required"))
	if err != nil {
		return nil, err
	}
	if !ok {
		obj.Required = false
	} else {
		obj.Required, ok = util.StarlarkAsBool(v)
		if !ok {
			return nil, fmt.Errorf("Field 'required' not a bool! Not a schema-object?")
		}
	}
	// default
	v, ok, err = dict.Get(starlark.String("default"))
	if err != nil {
		return nil, err
	}
	if !ok {
		obj.Default, err = getDefaultValue(obj.Type)
		if err != nil {
			return nil, err
		}
	} else {
		// TODO
	}
	return obj, nil
}

func (obj *SchemaObject) ToDict() starlark.Value {
	dict := starlark.NewDict(16)
	dict.SetKey(starlark.String("type"), starlark.String(obj.Type))
	dict.SetKey(starlark.String("required"), starlark.Bool(obj.Required))
	if obj.Default != nil {
		// TODO: implement helper function to convert any type to starlark value!
	}
	return dict
}

func (obj *SchemaObject) UnpackFromArgs(args starlark.Tuple, kwargs []starlark.Tuple) error {
	if err := starlark.UnpackArgs(obj.Type, args, kwargs, "required?", &obj.Required, "default?", &obj.Default); err != nil {
		return err
	}
	return nil
}

func primitiveSchemaObjectFunc(name string, defaultValue interface{}) func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

		obj := NewSchemaObject(name, defaultValue)
		if err := obj.UnpackFromArgs(args, kwargs); err != nil {
			return nil, err
		}
		return obj.ToDict(), nil
	}
}

// func schemaList(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
// func schemaDict(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {

func getDefaultValue(name string) (interface{}, error) {
	switch name {
	case "string", "filename":
		return defaultString, nil
	case "bool":
		return defaultBool, nil
	case "float":
		return defaultFloat, nil
	case "int":
		return defaultInt, nil
	default:
		return "", fmt.Errorf("Unknown type: %s", name)
	}
}

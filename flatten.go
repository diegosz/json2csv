package main

import (
	"log"
	"reflect"
	"sort"
	"strconv"

	"json2csv/jsonpointer"
)

type mapKeys []reflect.Value

func (k mapKeys) Len() int           { return len(k) }
func (k mapKeys) Swap(i, j int)      { k[i], k[j] = k[j], k[i] }
func (k mapKeys) Less(i, j int) bool { return k[i].String() < k[j].String() }

func sortedMapKeys(v reflect.Value) []reflect.Value {
	var keys mapKeys = v.MapKeys()
	sort.Sort(keys)
	return keys
}

type keyValue map[string]interface{}

func (kv keyValue) Keys() []string {
	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	return keys
}

func flatten(obj interface{}) keyValue {
	f := make(keyValue, 0)
	key := jsonpointer.JSONPointer{}
	_flatten(f, obj, key)
	return f
}

func _flatten(out keyValue, obj interface{}, key jsonpointer.JSONPointer) {
	value, ok := obj.(reflect.Value)
	if !ok {
		value = reflect.ValueOf(obj)
	}
	for value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Map:
		_flattenMap(out, value, key)
	case reflect.Slice:
		_flattenSlice(out, value, key)
	case reflect.String:
		out[key.String()] = value.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out[key.String()] = value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out[key.String()] = value.Uint()
	case reflect.Float32, reflect.Float64:
		out[key.String()] = value.Float()
	case reflect.Bool:
		out[key.String()] = value.Bool()
	default:
		log.Fatal("Unknown kind: ", value.Kind())
	}
}

func _flattenMap(out map[string]interface{}, value reflect.Value, prefix jsonpointer.JSONPointer) {
	keys := sortedMapKeys(value)
	for _, key := range keys {
		pointer := prefix.Clone()
		pointer.AppendString(key.String())
		_flatten(out, value.MapIndex(key).Interface(), pointer)
	}
}

func _flattenSlice(out map[string]interface{}, value reflect.Value, prefix jsonpointer.JSONPointer) {
	for i := 0; i < value.Len(); i++ {
		pointer := prefix.Clone()
		pointer.AppendString(strconv.Itoa(i))
		_flatten(out, value.Index(i).Interface(), pointer)
	}
}
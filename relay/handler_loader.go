package main

import (
	"fmt"
	"os"
	"plugin"
	"reflect"
	"strings"

	"farseer/handlers"
	protos "farseer/protos"

	"github.com/charmbracelet/log"
)

func LoadHandler(name string, messages chan *protos.GossipMessage, ll log.Logger) error {
	pl, err := plugin.Open(fmt.Sprintf("../compiled_handlers/%s.so", name))
	if err != nil {
		return err
	}

	plEventHandlersSymbol, err := pl.Lookup("PluginHandler")
	if err != nil {
		log.Error("Couldn't find the symbol containing the event handlers!", "PluginName", name)
		return err
	}

	plEventHandlers := *plEventHandlersSymbol.(*handlers.Handler)

	go plEventHandlers.HandleMessages(messages, ll)

	return nil
}

func LoadHandlersFromConf(conf Config, messages chan *protos.GossipMessage, ll log.Logger) {
	keys := make([]string, 0, len(conf.Servers))
	for k := range conf.Servers {
		isKEnabled := conf.Servers[k].(map[string]interface{})["enabled"]
		if isKEnabled == true && reflect.TypeOf(isKEnabled).String() == "string" {
			keys = append(keys, isKEnabled.(string))
		}
	}

	availableHandlers, err := ListCompiledHandlers()
	if err != nil {
		ll.Error("Couldn't get available handlers from folder!")
		return
	}

	for _, el := range intesectionOfArrays(keys, availableHandlers) {
		err = LoadHandler(el, messages, ll)
		if err != nil {
			ll.Error("Couldn't load handlers from conf! |", "Error", err)
			return
		}
	}
}

func intesectionOfArrays[T comparable](a []T, b []T) []T {
	set := make([]T, 0)

	for _, v := range a {
		if containsGeneric(b, v) {
			set = append(set, v)
		}
	}

	return set
}

func containsGeneric[T comparable](b []T, e T) bool {
	for _, v := range b {
		if v == e {
			return true
		}
	}
	return false
}

func ListCompiledHandlers() ([]string, error) {
	plList := []string{}

	dirEntries, err := os.ReadDir("../compiled_handlers")
	if err != nil {
		return []string{}, err
	}

	for _, entry := range dirEntries {
		entryI, err := entry.Info()
		if err != nil {
			return []string{}, err
		}
		plList = append(plList, strings.Split(entryI.Name(), ".")[0])
	}

	return plList, nil
}

package main

import (
	"fmt"
	"os"
	"plugin"
	"strings"

	"farseer/config"
	"farseer/handlers"
	"farseer/utils"
	protos "farseer/protos"

	"github.com/charmbracelet/log"
)

func LoadHandlersFromConf(conf config.Config, messages chan *protos.GossipMessage, ll log.Logger) error {
	keys := conf.GetHandlers()

	availableHandlers, err := ListCompiledHandlers()
	if err != nil {
		ll.Error("Couldn't get available handlers from folder!")
		return err
	}

	ll.Debug("Available handlers! |", "Handlers", availableHandlers)

	for _, el := range utils.IntersectionOfArrays(keys, availableHandlers) {
		ll.Debug("Loading handlers! |", "Element", el)
		err = LoadHandler(el, messages, ll, conf)
		if err != nil {
			ll.Error("Couldn't load handlers from conf! |", "Error", err)
			return err
		}
	}

	return nil
}

func LoadHandler(name string, messages chan *protos.GossipMessage, ll log.Logger, conf config.Config) error {
	pl, err := plugin.Open(fmt.Sprintf("../compiled_handlers/%s.so", name))
	if err != nil {
		return err
	}

	ll.Debug("Opening shared lib! |", "Name", name, "Handlers", conf.GetHandlers())

	plEventHandlersSymbol, err := pl.Lookup("PluginHandler")
	if err != nil {
		log.Error("Couldn't find the symbol containing the event handlers!", "PluginName", name)
		return err
	}

	plEventHandlers := *plEventHandlersSymbol.(*handlers.Handler)

	params := conf.GetParams(name)
	if params == nil {
		params = map[string]interface{}{}
	}
	go plEventHandlers.HandleMessages(messages, ll, params)

	return nil
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

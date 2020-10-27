package main

type announcement struct {
	Author             string `json:"author"`
	Message            string `json:"message"`
	PluginInternalName string `json:"pluginInternalName"`
}

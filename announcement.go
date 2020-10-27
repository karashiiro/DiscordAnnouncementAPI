package main

type announcement struct {
	Author             *author `json:"author"`
	Message            string  `json:"message"`
	PluginInternalName string  `json:"pluginInternalName"`
	Link               string  `json:"link"`
}

type author struct {
	Nickname  string `json:"nickname"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl"`
}

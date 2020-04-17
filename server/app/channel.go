package app

type Channel struct {
	Name      string `json:"name"`
	Creator   string `json:"creator"`
	Clients   []*Client
	ClientMap map[string]*Client
}

func (channel *Channel) addClient(client *Client) {
	channel.Clients = append(channel.Clients, client)
	channel.ClientMap[client.Addr.String()] = client
}

func (channel Channel) getClientsNames() []string {
	var names []string
	for _, client := range channel.Clients {
		names = append(names, client.Name)
	}
	return names
}

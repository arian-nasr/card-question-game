package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocketMessage defines the structure for our JSON messages
type WebSocketMessage struct {
	Event   string      `json:"event"`   // The name of the event (e.g., "joinLobby", "chatMessage")
	Payload interface{} `json:"payload"` // The actual data for the event (can be any structure)
}

var connManager = NewConnectionManager()
var lobbyManager = NewLobbyManager()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ConnectionManager struct {
	Clients map[string]*Client
}

type Client struct {
	Conn  *websocket.Conn
	ID    string
	Name  string
	Lobby *Lobby
}

type Lobby struct {
	ID      string
	Clients map[string]*Client
}

type LobbyManager struct {
	Lobbies map[string]*Lobby
}

func CreateLobby(id string) *Lobby {
	return &Lobby{
		ID:      id,
		Clients: make(map[string]*Client),
	}
}

func NewLobbyManager() *LobbyManager {
	return &LobbyManager{
		Lobbies: make(map[string]*Lobby),
	}
}

func (lm *LobbyManager) AddLobby(lobby *Lobby) {
	lm.Lobbies[lobby.ID] = lobby
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		Clients: make(map[string]*Client),
	}
}

func (cm *ConnectionManager) AddClient(client *Client) {
	cm.Clients[client.ID] = client
}

func (cm *ConnectionManager) RemoveClient(client *Client) {
	delete(cm.Clients, client.ID)
}

func (cm *ConnectionManager) SendToClient(clientID string, message []byte) {
	client, exists := cm.Clients[clientID]
	if exists {
		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("Error while writing message to client:", err)
		}
	} else {
		fmt.Println("Client not found:", clientID)
	}
}

//func NewLobby() *Lobby {
//	return &Lobby{
//		Clients: make(map[string]*Client),
//	}
//}

func (l *Lobby) AddClient(client *Client) {
	l.Clients[client.ID] = client
	client.Lobby = l
}

func (l *Lobby) RemoveClient(client *Client) {
	delete(l.Clients, client.ID)
	client.Lobby = nil
}

func (l *Lobby) SendToLobby(message []byte) {
	for _, client := range l.Clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("Error while writing message to client:", err)
		}
	}
}

/* func (l *Lobby) GetRandomQuestion() *Question {
	i := rand.Intn(len(l.Questions))
	q := l.Questions[i]
	l.Questions = slices.Delete(l.Questions, i, i+1)

	return q
}

func (l *Lobby) LoadQuestionsFromFile() error {
	lines, err := readLines("./questions.txt")
	if err != nil {
		return err
	}

	for _, line := range lines {
		q := CreateQuestion(line)
		l.AddQuestion(q)
	}

	return nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
} */

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	clientID := fmt.Sprintf("client-%p", conn)
	client := &Client{
		Conn: conn,
		ID:   clientID,
	}

	connManager.AddClient(client)
	defer func() {
		conn.Close()
		// don't forget to remove the client from the lobby too
		connManager.RemoveClient(client)
	}()

	welcome := WebSocketMessage{
		Event:   "welcome",
		Payload: fmt.Sprintf("Your ID is: %s", clientID),
	}

	err = conn.WriteJSON(welcome)
	if err != nil {
		fmt.Println("Error while writing welcome message:", err)
		return
	}

	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}

		switch msg.Event {

		case "hello":
			fmt.Println("Hello message received:", msg.Payload)

		case "joinLobby":
			fmt.Println(client.ID)
			lobbyCode := msg.Payload.(string)
			lobby, exists := lobbyManager.Lobbies[lobbyCode]
			if !exists {
				fmt.Println("Lobby not found:", lobbyCode)
				break
			}
			lobby.AddClient(client)
			// send response with payload of all players in the lobby
			players := make([]string, 0, len(lobby.Clients))
			for _, c := range lobby.Clients {
				players = append(players, c.ID)
			}
			lobbyResponse := WebSocketMessage{
				Event:   "joinLobby",
				Payload: players,
			}
			err = conn.WriteJSON(lobbyResponse)
			if err != nil {
				fmt.Println("Error while writing lobby response:", err)
			}

		default:
			fmt.Println("Unknown event:", msg.Event)

		}

		err = conn.WriteJSON(msg)
		if err != nil {
			fmt.Println("Error while writing message:", err)
			break
		}
	}

	/* for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			break
		}
		err = conn.WriteMessage(messageType, msg)
		if err != nil {
			fmt.Println("Error while writing message:", err)
			break
		}
	} */
}

func main() {
	// create test lobby with ID "test123"
	lobby := CreateLobby("test123")
	lobbyManager.AddLobby(lobby)
	http.HandleFunc("/ws", handleWebSocket)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

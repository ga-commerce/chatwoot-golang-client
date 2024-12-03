package chatwootclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Please note that certain functions like to add labels or assign agents are blocked when using an Agent Bot Token
// therefore an AgentToken has to be provided. The client uses the AgentBotToken wherever possible.
type ChatwootClient struct {
	BaseUrl string
}

func NewChatwootClient(baseUrl string) ChatwootClient {
	return ChatwootClient{
		baseUrl,
	}
}

func NewChatwootClientWithAgentToken(baseUrl string) ChatwootClient {
	return ChatwootClient{
		baseUrl,
	}
}

type CreateContactRequest struct {
	InboxID          int         `json:"inbox_id"`
	Name             string      `json:"name,omitempty"`
	EMail            string      `json:"email,omitempty"`
	PhoneNumber      string      `json:"phone_number,omitempty"`
	Avatar           string      `json:"avatar,omitempty"`
	AvatarUrl        string      `json:"avatar_url,omitempty"`
	Identifier       string      `json:"identifier,omitempty"`
	CustomAttributes interface{} `json:"custom_attributes,omitempty"`
}

type CreateContactResponse struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	Contact Contact `json:"contact"`
}

type Contact struct {
	ID             int            `json:"id"`
	ContactInboxes []ContactInbox `json:"contact_inboxes"`
}

type ContactInbox struct {
	SourceID string `json:"source_id"`
}

func (client *ChatwootClient) CreateContact(accountId int64, agentToken string, createContactRequest CreateContactRequest) (CreateContactResponse, error) {

	url := fmt.Sprintf("%s/api/v1/accounts/%v/contacts", client.BaseUrl, accountId)

	requestJSON, err := json.Marshal(createContactRequest)

	if err != nil {
		return CreateContactResponse{}, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestJSON))

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("api_access_token", agentToken)

	if err != nil {
		return CreateContactResponse{}, err
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return CreateContactResponse{}, err
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return CreateContactResponse{}, err
	}

	var createContactResponse CreateContactResponse

	if err := json.Unmarshal(body, &createContactResponse); err != nil {
		return CreateContactResponse{}, err
	}

	return createContactResponse, nil

}

type CreateNewConversationRequest struct {
	SourceID  string `json:"source_id"`
	InboxID   int    `json:"inbox_id"`
	ContactID string `json:"contact_id,omitempty"`
	Status    string `json:"status,omitempty"`
}

type CreateNewConversationResponse struct {
	ID        int `json:"id"`
	AccountId int `json:"account_id"`
	InboxId   int `json:"inbox_id"`
}

func (client *ChatwootClient) CreateNewConversation(accountId int64, agentBotToken string, createNewConversationRequest CreateNewConversationRequest) (CreateNewConversationResponse, error) {

	url := fmt.Sprintf("%s/api/v1/accounts/%v/conversations", client.BaseUrl, accountId)

	requestJSON, err := json.Marshal(createNewConversationRequest)

	if err != nil {
		return CreateNewConversationResponse{}, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestJSON))

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("api_access_token", agentBotToken)

	if err != nil {
		return CreateNewConversationResponse{}, err
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return CreateNewConversationResponse{}, err
	}

	if response.StatusCode != 200 {
		return CreateNewConversationResponse{}, errors.New("Request failed" + response.Status)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return CreateNewConversationResponse{}, err
	}

	var createNewConversationResponse CreateNewConversationResponse

	if err := json.Unmarshal(body, &createNewConversationResponse); err != nil {
		return CreateNewConversationResponse{}, err
	}

	return createNewConversationResponse, nil

}

type GetMessagesResponse struct {
	Meta    interface{}      `json:"meta"`
	Payload ChatwootMessages `json:"payload"`
}

type ChatwootMessages []struct {
	Id          int         `json:"id"`
	Content     string      `json:"content"`
	ContentType string      `json:"content_type,omitempty"`
	Private     bool        `json:"private,omitempty"`
	Sender      interface{} `json:"sender,omitempty"`
}

func (client *ChatwootClient) GetMessages(accountId int64, conversationId int64, agentToken string) (ChatwootMessages, error) {

	url := fmt.Sprintf("%s/api/v1/accounts/%v/conversations/%v/messages", client.BaseUrl, accountId, conversationId)

	request, _ := http.NewRequest(http.MethodGet, url, nil)

	request.Header.Add("api_access_token", agentToken)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New("Request failed" + response.Status)
	}

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var getMessagesResponse GetMessagesResponse

	if err := json.Unmarshal(responseBody, &getMessagesResponse); err != nil {
		return nil, err
	}

	return getMessagesResponse.Payload, nil

}

// Struct that allows to build minimal create message requests.
type CreateNewMessageRequest struct {
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
	Private     bool   `json:"private"`
}

type CreateNewMessageResponse struct {
	ID          int    `json:"id"`
	Content     string `json:"content"`
	MessageType int    `json:"message_type"` // Chatwoot 2.17.1 returns integers as message type in contrast to the API documentation
	Private     bool   `json:"private"`
}

func NewCreateNewMessageRequest(content string, messageType string, private bool) CreateNewMessageRequest {
	return CreateNewMessageRequest{
		Content:     content,
		MessageType: messageType,
		Private:     private,
	}
}

func (client *ChatwootClient) CreateNewMessage(accountId int64, conversationId int64, agentBotToken string, createMessageRequest CreateNewMessageRequest) (CreateNewMessageResponse, error) {

	url := fmt.Sprintf("%s/api/v1/accounts/%v/conversations/%v/messages", client.BaseUrl, accountId, conversationId)

	requestBodyJSON, err := json.Marshal(createMessageRequest)

	if err != nil {
		return CreateNewMessageResponse{}, err
	}

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyJSON))

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("api_access_token", agentBotToken)

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return CreateNewMessageResponse{}, err
	}

	if response.StatusCode != 200 {
		return CreateNewMessageResponse{}, errors.New("Request failed" + response.Status)
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return CreateNewMessageResponse{}, err
	}

	var createNewMessageResponse CreateNewMessageResponse

	if err := json.Unmarshal(body, &createNewMessageResponse); err != nil {
		return CreateNewMessageResponse{}, err
	}

	return createNewMessageResponse, nil

}

func (client *ChatwootClient) CreateOutgoingMessage(accountId int64, conversationId int64, agentBotToken string, content string) (CreateNewMessageResponse, error) {

	return client.CreateNewMessage(accountId, conversationId, agentBotToken, NewCreateNewMessageRequest(
		content,
		"outgoing",
		false,
	))

}

func (client *ChatwootClient) CreateOutgoingPrivateMessage(accountId int64, conversationId int64, agentBotToken string, content string) (CreateNewMessageResponse, error) {

	return client.CreateNewMessage(accountId, conversationId, agentBotToken, NewCreateNewMessageRequest(
		content,
		"outgoing",
		true,
	))

}

func (client *ChatwootClient) CreateIncomingMessage(accountId int64, conversationId int64, agentBotToken string, content string) (CreateNewMessageResponse, error) {

	return client.CreateNewMessage(accountId, conversationId, agentBotToken, NewCreateNewMessageRequest(
		content,
		"incoming",
		false,
	))

}

func (client *ChatwootClient) CreateIncomingPrivateMessage(accountId int64, conversationId int64, agentBotToken string, content string) (CreateNewMessageResponse, error) {

	return client.CreateNewMessage(accountId, conversationId, agentBotToken, NewCreateNewMessageRequest(
		content,
		"incoming",
		true,
	))

}

type AddLabelsRequest struct {
	Labels []string `json:"labels"`
}

func (client *ChatwootClient) AddLabels(accountId int64, conversationId int64, agentToken string, labels []string) error {

	if agentToken == "" {
		return errors.New("agentToken is empty. Adding labels requires a Chatwoot agent token")
	}

	url := fmt.Sprintf("%s/api/v1/accounts/%v/conversations/%v/labels", client.BaseUrl, accountId, conversationId)

	requestBody := AddLabelsRequest{
		Labels: labels,
	}

	requestBodyJSON, err := json.Marshal(requestBody)

	if err != nil {
		return err
	}

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyJSON))

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("api_access_token", agentToken)

	response, err := http.DefaultClient.Do(request)

	if response.StatusCode != 200 {
		return errors.New("Request failed" + response.Status)
	}

	return err

}

func (client *ChatwootClient) AddLabel(accountId int64, conversationId int64, agentToken string, label string) error {

	if agentToken == "" {
		return errors.New("agentToken is empty. Adding labels requires a Chatwoot agent token")
	}

	url := fmt.Sprintf("%s/api/v1/accounts/%v/conversations/%v/labels", client.BaseUrl, accountId, conversationId)

	requestBody := AddLabelsRequest{
		Labels: []string{label},
	}

	requestBodyJSON, err := json.Marshal(requestBody)

	if err != nil {
		return err
	}

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyJSON))

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("api_access_token", agentToken)

	response, err := http.DefaultClient.Do(request)

	if response.StatusCode != 200 {
		return errors.New("Request failed" + response.Status)
	}

	return err

}

func (client *ChatwootClient) Assign(accountId int64, conversationId int64, agentToken string, assignee_id int) error {

	if agentToken == "" {
		return errors.New("agentToken is empty. Adding assignments requires a Chatwoot agent token")
	}

	url := fmt.Sprintf("%s/api/v1/accounts/%v/conversations/%v/assignments", client.BaseUrl, accountId, conversationId)

	requestBody := fmt.Sprintf(`{"assignee_id": %v}`, assignee_id)

	requestBodyAsBytes := []byte(requestBody)

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyAsBytes))

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("api_access_token", agentToken)

	response, err := http.DefaultClient.Do(request)

	if response.StatusCode != 200 {
		return errors.New("Request failed" + response.Status)
	}

	return err

}

func (client *ChatwootClient) AssignTeam(accountId int64, conversationId int64, agentToken string, team_id int) error {

	if agentToken == "" {
		return errors.New("agentToken is empty. Adding assignments requires a Chatwoot agent token")
	}

	url := fmt.Sprintf("%s/api/v1/accounts/%v/conversations/%v/assignments", client.BaseUrl, accountId, conversationId)

	requestBody := fmt.Sprintf(`{"team_id": %v}`, team_id)

	requestBodyAsBytes := []byte(requestBody)

	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBodyAsBytes))

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("api_access_token", agentToken)

	response, err := http.DefaultClient.Do(request)

	if response.StatusCode != 200 {
		return errors.New("Request failed" + response.Status)
	}

	return err

}

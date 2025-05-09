package chatwootclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateContact(t *testing.T) {

	// configure mocked chatwoot server

	createContactResponse := CreateContactResponse{
		Payload: Payload{
			Contact: Contact{
				ContactInboxes: []ContactInbox{
					{
						SourceID: "42",
					},
				},
			},
		},
	}

	body, _ := json.Marshal(createContactResponse)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write(body)

	}))

	defer server.Close()

	// client set up

	client := ChatwootClient{
		BaseUrl: server.URL,
	}

	response, err := client.CreateContact(1, "", CreateContactRequest{
		InboxID: 1,
		Name:    "Unit Test Contact",
	})

	if err != nil {
		t.FailNow()
	}

	if response.Payload.Contact.ContactInboxes[len(response.Payload.Contact.ContactInboxes)-1].SourceID != "42" {
		t.FailNow()
	}

}

func TestSendImageMessage(t *testing.T) {
	client := ChatwootClient{
		BaseUrl: "xxxx.com",
	}

	res, err := client.SendImageMessage(2, 16, "1kA6UJ2Mgeu8hmjMWLy3edPx", "https://cdn.shopify.com/s/files/1/0723/2808/5795/files/20241129170255.jpg?v=1732875275", "Silk dress $468❤️\nWould you like one?\nIt’s one of our hot-selling items😆")
	if err != nil {
		fmt.Println("message err: ", err)
		return
	}
	fmt.Println("res: ", res)
}

func TestSendTextMessage(t *testing.T) {
	client := ChatwootClient{
		BaseUrl: "https://xxx.xxx.com",
	}

	res, err := client.CreateNewMessage(2, 14, "xxx", CreateNewMessageRequest{
		Content:     "test content",
		MessageType: "outgoing",
		Private:     false,
	})
	if err != nil {
		fmt.Println("message err: ", err)
		return
	}
	fmt.Println("res: ", res)
}

func TestSendNotification(t *testing.T) {
	client := ChatwootClient{
		BaseUrl: "http://localhost:3000",
	}

	err := client.SendNotification(3, "bHC2BQpqncp7K9JRx7ua1CF2", SendNotificationRequest{
		UserId:           2,
		NotificationType: "insufficient_gpt_balance",
		PrimaryActorType: "User",
		PrimaryActorId:   3,
	})
	if err != nil {
		fmt.Println("message err: ", err)
		return
	}
}

func TestSendTips(t *testing.T) {
	client := ChatwootClient{
		BaseUrl: "http://localhost:3000",
	}

	err := client.SendConversationTips(3, 2, "bHC2BQpqncp7K9JRx7ua1CF2", "enabled")
	if err != nil {
		fmt.Println("message err: ", err)
		return
	}
}

func TestUpdateConversationAIDisabled(t *testing.T) {
	client := ChatwootClient{
		BaseUrl: "http://localhost:3000",
	}

	err := client.UpdateConversationAIDisabled(1, 1, "zFNvsxjwzvJFG5Z3hTezcgPt", false)
	if err != nil {
		fmt.Println("message err: ", err)
		return
	}
}

package chat

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"github.com/tjons/text-to-trade/pkg/ai"
	"github.com/tjons/text-to-trade/pkg/api/chat"
	"github.com/tjons/text-to-trade/pkg/gen"
	"github.com/tjons/text-to-trade/pkg/model"
	"gorm.io/gorm"
)

type chatHandler struct {
	chat.UnimplementedChatServer
	db *gorm.DB
}

func NewChatServer(db *gorm.DB) chat.ChatServer {
	return &chatHandler{db: db}
}

func (s *chatHandler) SendChatMessage(ctx context.Context, req *chat.Question) (*chat.Answer, error) {
	g, err := ai.Gemini(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	aiModel := g.GenerativeModel("gemini-1.0-pro")
	session := aiModel.StartChat()
	wakeupResponse, err := session.SendMessage(ctx, genai.Text(
		`You are an AI chatbot. 
		You are here to help me convert user text into API actions. 
		I will provide you with an OpenAPI spec and a user's text.
		You should respond with the appropriate API route, HTTP method, and HTTP body formatted as JSON like this:
		{
			"route": "$ROUTE",
			"method": "$METHOD",
			"body": "$BODY"
		}
		Do NOT include any code formatting or comments in your response. Just return the response as JSON.
		If you understand these instructions, please respond with "ack".`,
	))
	if err != nil {
		log.Printf("Error waking up chatbot: %v", err)
		return nil, err
	}

	wakeupResponseText, ok := wakeupResponse.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		log.Printf("Unexpected response type")
		return nil, fmt.Errorf("Unexpected response type")
	}

	if !strings.Contains(string(wakeupResponseText), "ack") {
		log.Printf("Chatbot failed to wake up: %v", wakeupResponse.Candidates[0].Content.Parts[0].(genai.Text))
		return nil, fmt.Errorf("Chatbot failed to wake up")
	}

	actualResponse, err := session.SendMessage(ctx, genai.Text(fmt.Sprintf(
		`
		OpenAPI spec: "%s"


		User's text: "%s".`,
		gen.WatchlistSwagger, req.Message,
	)))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	respText, ok := actualResponse.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		log.Printf("Unexpected response type")
		return nil, fmt.Errorf("Unexpected response type")
	}

	defer func() {
		chatHistory := model.ChatMessage{
			Question: req.Message,
			Answer:   string(respText),
		}
		if err := s.db.Create(&chatHistory).Error; err != nil {
			log.Println(err)
		}
	}()

	return &chat.Answer{Message: string(respText)}, nil
}

func (s *chatHandler) SendAdviceMessage(ctx context.Context, req *chat.Question) (*chat.Answer, error) {
	user := &model.User{}
	if err := s.db.First(user, req.UserId).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	g, err := ai.Gemini(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	aiModel := g.GenerativeModel("gemini-1.0-pro")
	session := aiModel.StartChat()
	wakeupResponse, err := session.SendMessage(ctx, genai.Text(
		`You are an expert investment advisor, with a strong track record of managing risk, outperforming the market, and understanding your clients.
		You will be provided with a user's text and asked to provide investment advice.
		Constrain your advice to each question to 300 words or less.
		If you understand these instructions, please respond with "ack".`,
	))
	if err != nil {
		log.Printf("Error waking up chatbot: %v", err)
		return nil, err
	}
	wakeupResponseText, ok := wakeupResponse.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		log.Printf("Unexpected response type")
		return nil, fmt.Errorf("Unexpected response type")
	}

	if strings.Contains(string(wakeupResponseText), "ack") {
		log.Printf("Chatbot failed to wake up: %v", wakeupResponse.Candidates[0].Content.Parts[0].(genai.Text))
		return nil, fmt.Errorf("Chatbot failed to wake up")
	}

	backgroundResponse, err := session.SendMessage(ctx, genai.Text(fmt.Sprintf(
		`Here's some context about the user's background: 
		The user is a %s investor with a %s risk tolerance and a %s timeframe.
		If you understand the user's background and are ready for their question, please response with "ack".`,
		user.Experience, user.Risk, user.Allocation,
	)))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	backgroundResponseText, ok := backgroundResponse.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		log.Printf("Unexpected response type")
		return nil, fmt.Errorf("Unexpected response type")
	}

	if strings.Contains(string(backgroundResponseText), "ack") {
		log.Printf("Chatbot failed to understand the user's background: %v", wakeupResponse.Candidates[0].Content.Parts[0].(genai.Text))
		return nil, fmt.Errorf("Chatbot failed to understand the user's background")
	}

	actualResponse, err := session.SendMessage(ctx, genai.Text(fmt.Sprintf(
		`Here is the users question: "%s".`,
		req.Message,
	)))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	actualResponseText, ok := actualResponse.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		log.Printf("Unexpected response type")
		return nil, fmt.Errorf("Unexpected response type")
	}

	return &chat.Answer{
		Message: string(actualResponseText),
	}, err
}

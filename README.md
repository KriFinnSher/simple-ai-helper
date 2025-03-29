# Customer Support Automation System - Backend Engineer Task

## Overview

This repository contains a backend solution to automate customer support using AI agents. The system is designed to handle customer queries in real-time, leverage natural language processing (NLP) to understand those queries, retrieve appropriate answers from a knowledge base, and escalate to a human agent when necessary.

The task includes designing the high-level architecture, implementing basic query handling, using an NLP model for intent recognition, and querying a mock database to provide responses. The system is designed to be scalable to handle multiple customer queries concurrently.

## Architecture Overview

The backend system consists of the following key components:

1. **AI Agent (NLP Model)**: 
   - The AI agent leverages a natural language processing model to analyze and understand customer queries. It is responsible for identifying the intent behind the query and deciding whether it can respond directly or needs escalation.

2. **Database**: 
   - A mock or in-memory database stores predefined responses to common queries. The AI agent queries this database to find relevant answers based on the identified intent. In a production system, this would be replaced with a real database (e.g., PostgreSQL, MySQL, or NoSQL solutions).

3. **Queueing System**: 
   - The system uses a queueing mechanism to handle multiple customer queries concurrently and ensures that the system can scale effectively as the number of incoming queries increases. This component ensures the system handles load efficiently and maintains performance.

4. **Scalability**: 
   - The system is designed to scale horizontally by adding more instances of the AI agent and database connections to handle increased traffic. Load balancing and proper handling of distributed systems are crucial in maintaining high availability and responsiveness.

## Key Requirements

- **AI/NLP for Query Understanding**: 
  - The system uses NLP techniques to understand the intent behind the customer query, enabling automated responses or escalating to a human agent when necessary.
  
- **Database Integration**: 
  - The system retrieves responses from a knowledge base (mock database) to provide relevant answers to customer queries.

- **Real-Time Processing**: 
  - The system handles customer queries in real-time and ensures timely responses.

- **Scalability**: 
  - The solution should efficiently handle concurrent queries to support a large volume of requests.

## Task Deliverables

1. **Architecture Diagram**: 
   - A high-level architecture diagram has already been created, depicting the flow of data and key components of the system.

2. **Code Implementation**: 
   - A code snippet simulating the query processing system, NLP-based intent recognition, and querying a mock database.

## Code Implementation

### 1. Query Handling (handler level)
```golang
func (h *Knowledge) Answer(ctx echo.Context) error {
	// processing request
	var req Request
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{Errors: "bad request body"})
	}

	// external API call
	intent, err := h.UseCase.ExtractIntent(ctx.Request().Context(), req.Message)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Errors: "failed to define intent"})
	}

	// database query
	answer, err := h.UseCase.GetAnswer(ctx.Request().Context(), intent)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{Errors: "failed to fetch answer"})
	}

	return ctx.JSON(http.StatusOK, Response{Answer: answer})
}
```
### 2. Intent extraction (usecase level)
```python
import spacy

nlp = spacy.load("en_core_web_md")

INTENTS = { # consider that we have some sort of FAQ in our app
    "pricing": ["cost", "price", "how much", "subscription fee"],
    "features": ["what can it do", "capabilities", "features", "functions"],
    "account": ["account issue", "login problem", "password reset"],
    "support": ["customer service", "help", "contact support", "assistance"],
    "refund": ["money back", "refund policy", "return payment", "cancel subscription"],
    "help": ["i need a human", "talk to agent", "talk to a real person"]
}

def extract_intent(message: str) -> str:
    doc = nlp(message.lower())
    best_intent = "help" # by default we escalate to a human agent
    best_score = 0.0

    for intent, examples in INTENTS.items():
        for example in examples:
            # checking similarity of user message and our predefined intents
            similarity = doc.similarity(nlp(example))
            if similarity > best_score:
                best_score = similarity
                best_intent = intent
    # we need at least 0.6 "similarity" score to chose particular intent
    return best_intent if best_score >= 0.6 else "help"


```

### 3. Database query (database level)
```golang
func (r *KnowledgeRepo) Get(ctx context.Context, intent string) (models.Knowledge, error) {
  // just getting predefined answer for our intent
	query, args, err := sq.Select("*").
		From("knowledge").
		Where(sq.Eq{"intent": intent}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return models.Knowledge{}, err
	}

	var knowledge models.Knowledge
	err = r.db.GetContext(ctx, &knowledge, query, args...)
	if err != nil {
		return models.Knowledge{}, err
	}

	return knowledge, nil
}
```

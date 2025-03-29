# Customer Support Automation System - Backend Engineer Task

## Overview

This repository contains a backend solution to automate customer support using AI agents. The system is designed to handle customer queries in real-time, leverage natural language processing (NLP) to understand those queries, retrieve appropriate answers from a knowledge base, and escalate to a human agent when necessary.

The task includes designing the high-level architecture, implementing basic query handling, using an NLP model for intent recognition, and querying a mock database to provide responses. The system is designed to be scalable to handle multiple customer queries concurrently.

## Architecture Overview

On the image bellow you can see flow of the main process:
- first, user sends message (request) to our app;
- secondly, we call to external API (which is my API on python) to use nlp model;
- thirdly, app need to check existance of received intent and fetch it from database;
- finaly we can send to user our response, which will be answer (message) to the user.

![flow drawio](https://github.com/user-attachments/assets/0dc09025-6f1f-4474-a89f-6bd41a924f9e)


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

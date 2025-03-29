import spacy

nlp = spacy.load("en_core_web_md")

INTENTS = {
    "pricing": ["cost", "price", "how much", "subscription fee"],
    "features": ["what can it do", "capabilities", "features", "functions"],
    "account": ["account issue", "login problem", "password reset"],
    "support": ["customer service", "help", "contact support", "assistance"],
    "refund": ["money back", "refund policy", "return payment", "cancel subscription"],
    "help": ["i need a human", "talk to agent", "talk to a real person"]
}

def extract_intent(message: str) -> str:
    doc = nlp(message.lower())
    best_intent = "help"
    best_score = 0.0

    for intent, examples in INTENTS.items():
        for example in examples:
            similarity = doc.similarity(nlp(example))
            if similarity > best_score:
                best_score = similarity
                best_intent = intent

    return best_intent if best_score >= 0.6 else "help"


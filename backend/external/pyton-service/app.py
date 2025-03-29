from fastapi import FastAPI
from pydantic import BaseModel
from nlp import extract_intent

app = FastAPI()

class Message(BaseModel):
    message: str

@app.post("/predict")
def predict(message: Message):
    intent = extract_intent(message.message)
    return {"intent": intent}

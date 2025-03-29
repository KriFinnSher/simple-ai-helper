CREATE TABLE knowledge (
                           intent TEXT PRIMARY KEY,
                           answer TEXT NOT NULL
);

CREATE INDEX idx_knowledge_intent ON knowledge(intent);

INSERT INTO knowledge (intent, answer) VALUES
                                           ('pricing', 'Our subscription starts at $9.99 per month.'),
                                           ('features', 'Our service includes AI-powered automation, analytics, and integrations.'),
                                           ('account', 'You can reset your password or update account details in the settings page.'),
                                           ('support', 'You can contact support via live chat or email us at support@example.com.'),
                                           ('refund', 'We offer a 30-day money-back guarantee on all subscriptions.'),
                                           ('help', 'I am transferring you to a human agent for further assistance.');

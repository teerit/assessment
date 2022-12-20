-- Sequence and defined type
CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);

INSERT INTO "expenses" ("id", "title", "amount", "note", "tags") VALUES (1, 'strawberry smoothie', 79.0, 'night market promotion discount 10 bath', ARRAY['food', 'beverage']);
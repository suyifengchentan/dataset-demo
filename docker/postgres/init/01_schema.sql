\connect visual_lab;

CREATE TABLE IF NOT EXISTS orders (
  id BIGSERIAL PRIMARY KEY,
  customer_name TEXT NOT NULL,
  status TEXT NOT NULL,
  tags TEXT[] NOT NULL DEFAULT '{}',
  payload JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO orders (customer_name, status, tags, payload)
VALUES
  ('Alice', 'paid', ARRAY['vip', 'east'], '{"channel":"web","amount":120.5,"items":[{"sku":"A-1","qty":2}]}'::jsonb),
  ('Bob', 'pending', ARRAY['north'], '{"channel":"app","amount":88.0,"items":[{"sku":"B-2","qty":1}]}'::jsonb),
  ('Carol', 'shipped', ARRAY['vip', 'south'], '{"channel":"store","amount":210.0,"items":[{"sku":"C-3","qty":4}]}'::jsonb)
ON CONFLICT DO NOTHING;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_publication WHERE pubname = 'visual_lab_pub') THEN
    CREATE PUBLICATION visual_lab_pub FOR TABLE orders;
  END IF;
END
$$;

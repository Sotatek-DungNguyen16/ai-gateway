#!/bin/bash
export AI_GATEWAY_API_KEY=your-api-key

echo '{"ai_model":"gemini-2.0-flash","ai_provider":"google","language":"javascript"}' > /tmp/meta.json
echo 'diff --git a/test.js b/test.js
+function unsafeEval(x) { 
+  return eval(x); 
+}' > /tmp/diff.txt

echo "Sending request..."
RESPONSE=$(curl -s -X POST http://localhost:8080/review \
  -H "X-API-Key: $AI_GATEWAY_API_KEY" \
  -F "metadata=$(cat /tmp/meta.json)" \
  -F "git_diff=@/tmp/diff.txt")

echo "=== RESPONSE ==="
echo "$RESPONSE" | head -100

echo ""
echo "=== DIAGNOSTIC COUNT ==="
echo "$RESPONSE" | jq '.diagnostics | length'

echo ""
echo "=== OVERVIEW ==="
echo "$RESPONSE" | jq -r '.overview'


#!/bin/bash
# Debug test to see actual response

# Create test files
cat > /tmp/meta.json << 'EOF'
{
  "ai_model": "gemini-2.0-flash",
  "ai_provider": "google",
  "language": "javascript"
}
EOF

cat > /tmp/diff.txt << 'EOF'
diff --git a/test.js b/test.js
+++ b/test.js
@@ -1,0 +1,3 @@
+function bad(x) {
+  return eval(x);
+}
EOF

echo "Sending request..."
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -X POST http://localhost:8080/review \
  -H "X-API-Key: $AI_GATEWAY_API_KEY" \
  -F "metadata=$(cat /tmp/meta.json)" \
  -F "git_diff=@/tmp/diff.txt")

echo ""
echo "=== RAW RESPONSE ==="
echo "$RESPONSE"
echo "=== END ==="

HTTP_STATUS=$(echo "$RESPONSE" | tail -n1 | sed 's/HTTP_STATUS://')
BODY=$(echo "$RESPONSE" | sed '$d')

echo ""
echo "Status: $HTTP_STATUS"
echo ""
echo "=== BODY ==="
echo "$BODY"
echo ""
echo "=== PRETTY ==="
echo "$BODY" | jq '.'


#!/bin/bash
# Test script for AI Gateway
export AI_GATEWAY_API_KEY=your-api-key
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo -e "${RED}❌ Error: jq is not installed${NC}"
    echo -e "${YELLOW}jq is required to parse JSON responses${NC}"
    echo ""
    echo "Install it with:"
    echo "  sudo apt-get install jq    # Debian/Ubuntu"
    echo "  brew install jq            # macOS"
    echo ""
    echo "Or view raw response:"
    echo "  ./debug-test.sh | less"
    exit 1
fi

echo -e "${YELLOW}Testing AI Gateway...${NC}\n"

# Check if server is running
echo -e "${YELLOW}1. Testing health endpoint...${NC}"
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
if echo "$HEALTH_RESPONSE" | grep -q "healthy"; then
    echo -e "${GREEN}✓ Health check passed${NC}"
else
    echo -e "${RED}✗ Health check failed${NC}"
    exit 1
fi

# Test review endpoint
echo -e "\n${YELLOW}2. Testing review endpoint...${NC}"

# Create test metadata
cat > /tmp/metadata.json << 'EOF'
{
  "ai_model": "gemini-2.0-flash",
  "ai_provider": "google",
  "language": "javascript",
  "review_mode": "file"
}
EOF

# Create comprehensive test diff covering all 6 categories
cat > /tmp/test.diff << 'EOF'
diff --git a/test.js b/test.js
index 1234567..abcdefg 100644
--- a/test.js
+++ b/test.js
@@ -1,3 +1,60 @@
+// ===== Category 1: POSSIBLE BUG =====
+// Issue: Null/undefined reference error
+function getUserEmail(user) {
+  return user.profile.email.toLowerCase(); // What if user, profile, or email is null?
+}
+
+// Issue: Array index out of bounds
+function getFirstItem(arr) {
+  return arr[0]; // What if array is empty?
+}
+
+// ===== Category 2: BEST PRACTICE =====
+// Issue: Poor naming convention
+function a(x, y) {
+  return x + y;
+}
+
+// Issue: Using var instead of const/let
+function calculate() {
+  var result = 0; // Should use const or let
+  return result;
+}
+
+// Issue: No error handling
+function parseJSON(str) {
+  return JSON.parse(str); // Should have try-catch
+}
+
+// ===== Category 3: ENHANCEMENT =====
+// Issue: Could use modern async/await
+function fetchUserData(id) {
+  return fetch(`/api/users/${id}`)
+    .then(res => res.json())
+    .then(data => data)
+    .catch(err => console.log(err));
+}
+
+// Issue: Could use optional chaining
+function getAddress(user) {
+  return user && user.profile && user.profile.address || 'N/A';
+}
+
+// ===== Category 4: POSSIBLE ISSUE =====
+// Issue: SQL Injection vulnerability
+function queryUser(username) {
+  return db.query("SELECT * FROM users WHERE name = '" + username + "'");
+}
+
+// Issue: Using eval (dangerous)
+function executeCode(code) {
+  return eval(code);
+}
+
+// ===== Category 5: MAINTAINABILITY =====
+// Issue: Magic numbers without explanation
+function calculateDiscount(price) {
+  if (price > 100) {
+    return price * 0.85; // What is 0.85? Why 100?
+  }
+  return price * 0.95; // What is 0.95?
+}
+
+// Issue: Long function with multiple responsibilities
+function processOrder(order) {
+  var total = order.items.reduce((sum, item) => sum + item.price, 0);
+  var tax = total * 0.08;
+  var shipping = total > 50 ? 0 : 5.99;
+  var discount = total > 100 ? total * 0.1 : 0;
+  var finalTotal = total + tax + shipping - discount;
+  sendEmail(order.user.email, finalTotal);
+  updateInventory(order.items);
+  logTransaction(order.id, finalTotal);
+  return finalTotal;
+}
+
+// Issue: Hard-coded configuration values
+function connectDatabase() {
+  return connect('localhost:3306', 'root', 'password123');
+}
+
+// ===== Category 6: PERFORMANCE =====
+// Issue: Inefficient nested loop (O(n²))
+function findDuplicates(arr1, arr2) {
+  let duplicates = [];
+  for (let i = 0; i < arr1.length; i++) {
+    for (let j = 0; j < arr2.length; j++) {
+      if (arr1[i] === arr2[j]) {
+        duplicates.push(arr1[i]);
+      }
+    }
+  }
+  return duplicates;
+}
+
+// Issue: Unnecessary re-renders / computations
+function processData(data) {
+  return data.map(item => {
+    const parsed = JSON.parse(item); // Parsing in every iteration
+    return parsed;
+  });
+}
+
+// Issue: Memory leak - not cleaning up event listeners
+function attachListener(element) {
+  element.addEventListener('click', function handler() {
+    console.log('clicked');
+  });
+  // No cleanup, listener stays in memory
+}
+
 function safeFunction() {
   return "Hello World";
 }
EOF

# Get API key from environment
if [ -z "$AI_GATEWAY_API_KEY" ]; then
    echo -e "${RED}Error: AI_GATEWAY_API_KEY environment variable not set${NC}"
    echo "Please set it with: export AI_GATEWAY_API_KEY=your-key"
    exit 1
fi

# Send review request
echo "Sending review request..."
REVIEW_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
  -X POST http://localhost:8080/review \
  -H "X-API-Key: $AI_GATEWAY_API_KEY" \
  -F "metadata=$(cat /tmp/metadata.json)" \
  -F "git_diff=@/tmp/test.diff")

# Extract status code
HTTP_STATUS=$(echo "$REVIEW_RESPONSE" | tail -n1 | sed 's/HTTP_STATUS://')
RESPONSE_BODY=$(echo "$REVIEW_RESPONSE" | sed '$d')

if [ "$HTTP_STATUS" = "200" ]; then
    echo -e "${GREEN}✓ Review request successful${NC}"
    
    # Debug: Save full response
    echo "$RESPONSE_BODY" > /tmp/test-response.json
    echo -e "${YELLOW}Debug: Full response saved to /tmp/test-response.json${NC}"
    
    # Show overview
    OVERVIEW=$(echo "$RESPONSE_BODY" | jq -r '.overview // "No overview provided"' 2>/dev/null)
    echo -e "\n${YELLOW}📋 Overview:${NC}"
    echo "$OVERVIEW"
    
    # Count total diagnostics
    DIAGNOSTIC_COUNT=$(echo "$RESPONSE_BODY" | jq '.diagnostics | length' 2>/dev/null || echo "0")
    echo -e "\n${GREEN}Found $DIAGNOSTIC_COUNT diagnostic(s)${NC}"
    
    # Show breakdown by category
    echo -e "\n${YELLOW}📊 Issues by Category:${NC}"
    for category in "possible-bug" "best-practice" "performance" "maintainability" "possible-issue" "enhancement"; do
        count=$(echo "$RESPONSE_BODY" | jq "[.diagnostics[] | select(.code.value == \"$category\")] | length" 2>/dev/null || echo "0")
        if [ "$count" -gt 0 ]; then
            # Format category name
            formatted=$(echo "$category" | tr '-' ' ' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)}1')
            echo -e "  ${GREEN}✓${NC} $formatted: ${YELLOW}$count${NC}"
        fi
    done
    
    # Show breakdown by severity
    echo -e "\n${YELLOW}🎯 Issues by Severity:${NC}"
    for severity in "ERROR" "WARNING" "INFO"; do
        count=$(echo "$RESPONSE_BODY" | jq "[.diagnostics[] | select(.severity == \"$severity\")] | length" 2>/dev/null || echo "0")
        if [ "$count" -gt 0 ]; then
            if [ "$severity" = "ERROR" ]; then
                color="${RED}"
            elif [ "$severity" = "WARNING" ]; then
                color="${YELLOW}"
            else
                color="${GREEN}"
            fi
            echo -e "  ${color}●${NC} $severity: ${color}$count${NC}"
        fi
    done
    
    # Show detailed issues
    echo -e "\n${YELLOW}📝 Detailed Issues:${NC}"
    echo "$RESPONSE_BODY" | jq -r '.diagnostics[] | "
File: \(.location.path)
Line: \(.location.range.start.line)
Category: \(.code.value)
Severity: \(.severity)
Message: \(.message)
" + (if .suggestion != "" and .suggestion != null then "Suggestion: \(.suggestion)\n" else "" end) + "---"' 2>/dev/null || echo "Unable to parse diagnostics"
    
else
    echo -e "${RED}✗ Review request failed with status $HTTP_STATUS${NC}"
    echo -e "${RED}Response:${NC}"
    echo "$RESPONSE_BODY"
    exit 1
fi

# Cleanup
rm -f /tmp/metadata.json /tmp/test.diff

echo -e "\n${GREEN}✅ All tests passed!${NC}"


# AI Code Review Categories

The AI Gateway reviews code against **6 specific categories** for every request. This ensures comprehensive, consistent feedback.

## 📋 Review Categories

### 1. 🐛 Possible Bug
**What it checks:**
- Null/undefined reference errors
- Array index out of bounds
- Division by zero
- Race conditions
- Off-by-one errors
- Unhandled edge cases
- Type mismatches
- Logic errors

**Example issues:**
```javascript
// ❌ Bad: No null check
function getUserEmail(user) {
  return user.profile.email.toLowerCase();
}

// ✅ Good: Safe access
function getUserEmail(user) {
  return user?.profile?.email?.toLowerCase() ?? 'N/A';
}
```

**Severity:** Usually `ERROR` or `WARNING`

---

### 2. ✨ Best Practice
**What it checks:**
- Naming conventions
- Code organization
- Design pattern misuse
- Language-specific idioms
- Error handling
- Use of modern syntax
- Deprecated API usage

**Example issues:**
```javascript
// ❌ Bad: Poor naming, no error handling
function a(x) {
  return JSON.parse(x);
}

// ✅ Good: Clear naming, error handling
function parseUserData(jsonString) {
  try {
    return JSON.parse(jsonString);
  } catch (error) {
    console.error('Failed to parse user data:', error);
    return null;
  }
}
```

**Severity:** Usually `WARNING` or `INFO`

---

### 3. 🚀 Performance
**What it checks:**
- Inefficient algorithms (O(n²), O(n³))
- Unnecessary loops
- Memory leaks
- N+1 query problems
- Blocking operations
- Large object copying
- Repeated computations
- Unnecessary re-renders

**Example issues:**
```javascript
// ❌ Bad: O(n²) complexity
function findDuplicates(arr1, arr2) {
  let duplicates = [];
  for (let i = 0; i < arr1.length; i++) {
    for (let j = 0; j < arr2.length; j++) {
      if (arr1[i] === arr2[j]) {
        duplicates.push(arr1[i]);
      }
    }
  }
  return duplicates;
}

// ✅ Good: O(n) with Set
function findDuplicates(arr1, arr2) {
  const set2 = new Set(arr2);
  return arr1.filter(item => set2.has(item));
}
```

**Severity:** Usually `WARNING` or `INFO`

---

### 4. 🔧 Maintainability
**What it checks:**
- Code complexity
- Magic numbers
- Hard-coded values
- Lack of documentation
- Unclear variable names
- Long functions
- Tight coupling
- Missing comments for complex logic

**Example issues:**
```javascript
// ❌ Bad: Magic numbers, no explanation
function calculatePrice(qty) {
  return qty * 19.99 * 1.08 + 5.99;
}

// ✅ Good: Named constants, clear intent
const ITEM_PRICE = 19.99;
const TAX_RATE = 0.08;
const SHIPPING_FEE = 5.99;

function calculatePrice(quantity) {
  const subtotal = quantity * ITEM_PRICE;
  const tax = subtotal * TAX_RATE;
  return subtotal + tax + SHIPPING_FEE;
}
```

**Severity:** Usually `WARNING` or `INFO`

---

### 5. ⚠️ Possible Issue
**What it checks:**
- Security vulnerabilities
- SQL injection risks
- XSS vulnerabilities
- CSRF issues
- Code smells
- Anti-patterns
- Potential future problems
- Resource leaks

**Example issues:**
```javascript
// ❌ Bad: SQL injection vulnerability
function queryUser(username) {
  return db.query(`SELECT * FROM users WHERE name = '${username}'`);
}

// ✅ Good: Parameterized query
function queryUser(username) {
  return db.query('SELECT * FROM users WHERE name = ?', [username]);
}
```

**Severity:** Usually `ERROR` or `WARNING`

---

### 6. 💡 Enhancement
**What it checks:**
- Optimization opportunities
- Better approaches available
- Modern syntax alternatives
- Missing features
- Code improvements
- Simplification opportunities

**Example issues:**
```javascript
// ❌ Could be improved: Promise chains
function fetchUserData(id) {
  return fetch(`/api/users/${id}`)
    .then(res => res.json())
    .then(data => data)
    .catch(err => console.log(err));
}

// ✅ Better: Async/await
async function fetchUserData(id) {
  try {
    const response = await fetch(`/api/users/${id}`);
    return await response.json();
  } catch (error) {
    console.error('Failed to fetch user:', error);
    throw error;
  }
}
```

**Severity:** Usually `INFO`

---

## 📊 How Issues Are Categorized

Each issue found will have:

1. **Category**: One of the 6 categories above
2. **Severity**: ERROR, WARNING, or INFO
3. **Line number**: Exact location in code
4. **Message**: Clear description of the issue
5. **Suggestion**: How to fix it (when applicable)

## 🎯 Severity Levels

### ERROR (🔴)
- Definite bugs
- Security vulnerabilities
- Critical performance issues
- Code that will cause runtime errors

### WARNING (🟡)
- Likely bugs
- Maintainability concerns
- Performance bottlenecks
- Anti-patterns
- Security concerns

### INFO (🟢)
- Best practice suggestions
- Enhancement opportunities
- Minor optimizations
- Code style improvements

## 🧪 Testing All Categories

The test script (`test-review.sh`) includes examples for all 6 categories:

```bash
# Make sure gateway is running
go run . &

# Set your API key
export AI_GATEWAY_API_KEY=your-key-here

# Run the comprehensive test
./test-review.sh
```

**Expected output:**
```
📊 Issues by Category:
  ✓ Possible Bug: 3
  ✓ Best Practice: 4
  ✓ Performance: 3
  ✓ Maintainability: 3
  ✓ Possible Issue: 2
  ✓ Enhancement: 2

🎯 Issues by Severity:
  ● ERROR: 5
  ● WARNING: 8
  ● INFO: 4
```

## 🔍 Language-Specific Focus

The AI adapts its review based on the detected language:

- **JavaScript/TypeScript**: Focuses on async/await, promises, React patterns
- **Python**: Focuses on list comprehensions, context managers, type hints
- **Java**: Focuses on OOP principles, exception handling, streams
- **Go**: Focuses on goroutines, error handling, interfaces
- **Others**: General best practices and common patterns

## 📝 Example Complete Review

```json
{
  "overview": "Found 17 issues across 6 categories. Critical security vulnerability in SQL query needs immediate attention.",
  "issues": [
    {
      "file": "test.js",
      "line": 3,
      "column": 10,
      "severity": "ERROR",
      "category": "possible-bug",
      "message": "Potential null pointer exception. user.profile.email may be null or undefined.",
      "suggestion": "Use optional chaining: user?.profile?.email?.toLowerCase() ?? 'N/A'"
    },
    {
      "file": "test.js",
      "line": 89,
      "column": 10,
      "severity": "ERROR",
      "category": "possible-issue",
      "message": "SQL injection vulnerability. User input is directly concatenated into SQL query.",
      "suggestion": "Use parameterized queries: db.query('SELECT * FROM users WHERE name = ?', [username])"
    }
    // ... more issues
  ]
}
```

## 🎓 Best Practices for Review

1. **Every category matters**: Don't focus only on bugs - maintainability and performance are equally important
2. **Actionable suggestions**: Each issue should have a clear fix
3. **Context-aware**: Consider the language and framework being used
4. **Prioritize by severity**: Fix ERRORs first, then WARNINGs, then INFOs
5. **Be constructive**: The goal is to improve code, not criticize

## 📚 Further Reading

- [System Prompt](internal/prompt/prompt.go#L15) - See how categories are defined
- [Test Examples](test-review.sh#L37) - Comprehensive test covering all categories
- [Integration Guide](../INTEGRATION_GUIDE.md) - How reviews work end-to-end

---

**Remember:** The AI reviews **ALL** code changes against **ALL 6 categories** for every request. This ensures nothing is missed! 🎯


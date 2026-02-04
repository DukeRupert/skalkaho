# Skalkaho SaaS Launch Readiness Analysis

**Date:** February 4, 2026
**Analyst:** Claude Code
**Product:** Skalkaho - Construction Quoting SaaS
**Current Status:** Feature-Complete MVP (Single-User)
**Target:** Production-Ready SaaS with Paying Customers

---

## Executive Summary

### Current Readiness Level: **35% Ready for Alpha Launch**

Skalkaho has a **solid, feature-rich MVP** with excellent domain logic for construction quoting. The core product works well and delivers real value. However, it is currently a **single-user application with no multi-tenancy, authentication, or billing infrastructure**.

**Critical Assessment:**
- **Product Core:** Strong (90% complete) - The quoting engine is sophisticated and well-implemented
- **SaaS Infrastructure:** Missing (10% complete) - No user accounts, tenant isolation, or payment processing
- **Security:** Basic (20% complete) - No authentication, CSRF protection, or rate limiting
- **Operational Readiness:** Moderate (40% complete) - Docker deployment exists but lacks monitoring, backups

**The Gap:** You're approximately **6-10 weeks of focused development** away from accepting your first paying customer. The good news: the hard problems (complex surcharge calculations, e-signatures, data modeling) are solved. The remaining work is mostly infrastructure plumbing.

### Top 5 Launch Blockers

| Priority | Blocker | Impact | Effort | Risk if Skipped |
|----------|---------|--------|--------|-----------------|
| 1 | Multi-tenancy | CRITICAL | XL (1-2 weeks) | Data leaks between customers |
| 2 | User Authentication | CRITICAL | L (4-5 days) | No access control at all |
| 3 | Stripe Billing | CRITICAL | L (4-5 days) | Cannot collect payments |
| 4 | PostgreSQL Migration | HIGH | M (2-3 days) | SQLite won't scale, deployment issues |
| 5 | Email Infrastructure | HIGH | M (2-3 days) | Can't send estimates or password resets |

**Total Estimated Effort for Alpha Launch:** 8-10 weeks (200-250 hours)

### Recommended Launch Strategy

**Option A: Rapid SaaS Launch (8-10 weeks)**
- Implement all 5 critical blockers
- Launch with manual onboarding
- Stripe integration for payment
- Basic email notifications
- **Target:** Self-serve signups with monthly billing

**Option B: Pilot-First Approach (4-6 weeks)**
- Skip multi-tenancy initially (single-tenant deployments)
- Add basic auth per instance
- Manual billing (invoice customers)
- Launch with 3-5 pilot customers
- **Target:** Learn and iterate before full SaaS build

**Recommendation:** Choose Option B if you want to validate pricing/market fit faster. Choose Option A if you're confident in the business model and want to scale immediately.

---

## Infrastructure Inventory

### 1. Multi-Tenancy

**Status:** ❌ **MISSING** - This is a single-user application

**Current State:**
- No `Organization` or `Tenant` model exists
- No tenant scoping on any data models
- All data is globally accessible
- Database has no `org_id` or `tenant_id` columns

**What This Means:**
Currently, if you deployed this as-is, every user would see every job, every client, every estimate. This is the #1 launch blocker.

**Gaps:**
- [ ] No tenant/organization data model
- [ ] No tenant context propagation through requests
- [ ] No tenant-scoped queries
- [ ] No tenant identification strategy (subdomain vs session)
- [ ] No tenant-aware migrations

**Launch Blocker:** ✅ **YES** - Absolutely cannot launch without this

**Estimated Effort:** **XL (1-2 weeks / 40-60 hours)**

**Implementation Notes:**
```
Required Changes:
1. Add organizations table (name, subdomain, stripe_customer_id, plan, status)
2. Add users table (email, password_hash, org_id, role)
3. Add org_id to ALL existing tables:
   - jobs, categories, line_items
   - clients, estimates, signature_requests
   - item_templates, price_imports
   - settings (one per org)
   - company_profile (one per org)
4. Update ALL queries to include WHERE org_id = ?
5. Add middleware to extract org_id from session
6. Add context propagation for tenant
7. Migration strategy for existing data
```

**Risk:** If you deploy without multi-tenancy and try to add it later, you'll need complex data migrations and risk downtime.

---

### 2. Authentication & Identity

**Status:** ❌ **MISSING** - No authentication exists

**Current State:**
- No user model
- No login/logout
- No password hashing
- No session management
- Routes are completely public

**What This Means:**
Anyone who knows the URL can access the application. Fine for MVP, catastrophic for production.

**Gaps:**
- [ ] User registration with email verification
- [ ] Secure password hashing (argon2id recommended)
- [ ] Login/logout with session cookies
- [ ] Password reset flow
- [ ] Session token security (HttpOnly, Secure, SameSite)
- [ ] Session storage (Redis or database)
- [ ] "Remember me" functionality

**Launch Blocker:** ✅ **YES** - Cannot launch without auth

**Estimated Effort:** **L (4-5 days / 32-40 hours)**

**Implementation Notes:**
```
Recommended Stack:
- Password hashing: golang.org/x/crypto/argon2
- Sessions: gorilla/sessions or custom JWT
- Email verification: Token-based (30-day expiry)

Routes to Protect:
- All /jobs/* routes
- All /estimates/* routes
- All /clients/* routes
- All /settings routes

Public Routes (Keep Open):
- /sign/{token} (e-signature public page)
- /health (monitoring)
- Marketing site (future)
```

**Acceptance Criteria:**
- User can register with email + password
- Email verification required before access
- Secure session cookies prevent hijacking
- Password reset via email works
- Session timeout after 30 days inactive

---

### 3. Billing & Subscriptions

**Status:** ❌ **MISSING** - No payment integration

**Current State:**
- No Stripe integration
- No subscription tracking
- No payment collection
- No access control based on plan status

**What This Means:**
You have no way to charge customers for the service you're providing.

**Gaps:**
- [ ] Stripe SDK integration
- [ ] Customer creation in Stripe
- [ ] Subscription creation and management
- [ ] Webhook endpoint for Stripe events
- [ ] Subscription status tracking (active, past_due, canceled)
- [ ] Trial period handling
- [ ] Dunning/grace period (failed payment retry logic)
- [ ] Subscription upgrade/downgrade flows
- [ ] Prorated billing logic

**Launch Blocker:** ✅ **YES** - For paid launch, absolutely

**Estimated Effort:** **L (4-5 days / 32-40 hours)**

**Alpha Launch Minimum:**
```
Must Have:
- Stripe customer creation on signup
- Single plan subscription ($19/mo)
- Webhook handling for subscription.updated
- Access gate: block app if subscription inactive
- Basic trial period (7-14 days)

Can Wait:
- Multiple plan tiers
- Upgrade/downgrade flows
- Proration logic
- Usage-based billing
- Invoice customization
```

**Implementation Notes:**
```
Subscription States to Handle:
- trialing: Allow full access
- active: Allow full access
- past_due: Grace period (3 days), show warning
- canceled: Block access, show reactivation prompt
- unpaid: Block access

Webhook Events (Critical):
- customer.subscription.created
- customer.subscription.updated
- customer.subscription.deleted
- invoice.payment_succeeded
- invoice.payment_failed
```

**Risk:** Webhook delivery can fail. You MUST implement idempotency keys and webhook signature verification.

---

### 4. Email Infrastructure

**Status:** ❌ **MISSING** - No email sending capability

**Current State:**
- No email service integration
- E-signature feature references "sending" but no actual email
- No transactional email templates

**What This Means:**
Critical user flows won't work:
- Can't verify email addresses
- Can't send password resets
- Can't deliver estimates to clients
- Can't notify of signature completions

**Gaps:**
- [ ] Email service provider integration (SendGrid, Postmark, AWS SES)
- [ ] Transactional email templates
- [ ] Email verification emails
- [ ] Password reset emails
- [ ] Estimate delivery emails
- [ ] Signature request emails
- [ ] Signature completion notifications
- [ ] Payment confirmation emails
- [ ] Failed payment warnings

**Launch Blocker:** ✅ **YES** - High priority

**Estimated Effort:** **M (2-3 days / 16-24 hours)**

**Recommended Service:**
```
Best Options for Small SaaS:
1. Postmark - Excellent deliverability, transactional focus
   Pricing: $15/mo for 10K emails

2. SendGrid - More features, slightly worse reputation
   Pricing: Free tier (100/day), then $15/mo

3. AWS SES - Cheapest, but requires more setup
   Pricing: $0.10 per 1000 emails

Recommendation: Postmark for simplicity and reliability
```

**Email Priorities:**
```
Critical (Launch Blockers):
1. Email verification (signup flow)
2. Password reset
3. Signature request delivery
4. Payment receipt

Important (Post-Launch):
5. Welcome email
6. Estimate sent notification
7. Signature completed notification
8. Trial ending reminder
9. Payment failed warning

Nice to Have:
- Weekly digest
- Feature announcements
- Onboarding sequences
```

**Templates Needed (Plain Text is Fine for Alpha):**
```
1. verify-email.txt
   Subject: Verify your Skalkaho account

2. password-reset.txt
   Subject: Reset your Skalkaho password

3. signature-request.txt
   Subject: Please review and sign your estimate

4. payment-receipt.txt
   Subject: Payment confirmation - Skalkaho subscription
```

---

### 5. Marketing Site & Public Pages

**Status:** 🟡 **PARTIAL** - No marketing site, some legal pages needed

**Current State:**
- Application UI exists and works well
- No landing page with pricing/signup
- No Terms of Service
- No Privacy Policy
- Public signature pages exist (good!)

**Gaps:**
- [ ] Landing page with value proposition
- [ ] Pricing page with tier comparison
- [ ] Signup page (with Stripe integration)
- [ ] Terms of Service
- [ ] Privacy Policy
- [ ] About page
- [ ] Help/documentation
- [ ] Route separation (public vs authenticated)

**Launch Blocker:** 🟡 **PARTIAL** - Legal pages are blockers, marketing site can be minimal

**Estimated Effort:** **M (2-3 days / 16-24 hours)**

**Alpha Launch Minimum:**
```
Must Have:
- Simple landing page (hero + pricing + CTA)
- Signup flow with Stripe
- Terms of Service (template + customization)
- Privacy Policy (template + customization)

Can Be Simple:
- One-page landing (no separate about/features)
- Generic legal docs (customize later)
- No blog or content marketing
- No help docs (email support only)
```

**Legal Documents:**
```
Resources:
- Use templates from Basecamp/Stripe legal docs
- Key sections:
  - Data ownership (customer owns their data)
  - Cancellation policy (cancel anytime)
  - Data retention (30 days after cancel)
  - No warranty (software "as is")
  - Limitation of liability
  - GDPR compliance (if EU customers)
```

**Recommendation:** Use a service like [Termly](https://termly.io) or [Iubenda](https://www.iubenda.com) to generate compliant legal docs quickly.

---

### 6. Onboarding Flow

**Status:** ❌ **MISSING** - No post-signup onboarding

**Current State:**
- Application assumes you know what to do
- No first-time user guidance
- No setup wizard

**Gaps:**
- [ ] Post-signup company profile setup
- [ ] First-time user tutorial/walkthrough
- [ ] Sample job/estimate (optional)
- [ ] Setup completion checklist
- [ ] Keyboard shortcuts help (exists but not onboarding)

**Launch Blocker:** ❌ **NO** - Can launch with manual onboarding

**Estimated Effort:** **S (1 day / 8 hours)** for minimal version

**Alpha Launch Strategy:**
```
Manual Onboarding Acceptable:
- Email new users after signup
- 15-minute onboarding call
- Share screen, walk through first quote
- Gather feedback on confusing parts

Automated (Post-Alpha):
- Interactive product tour
- Embedded help videos
- Contextual tooltips
- Sample data pre-populated
```

**Low-Effort Quick Win:**
```
Add to first login:
1. Welcome modal
2. "Create your first quote" button
3. Link to keyboard shortcuts (? key)
4. Support email prominently displayed
```

---

### 7. Security Hardening

**Status:** 🟡 **PARTIAL** - Basic logging exists, security features missing

**Current State:**
- ✅ Request ID tracking
- ✅ Structured logging
- ✅ Panic recovery
- ❌ No CSRF protection
- ❌ No rate limiting
- ❌ No security headers
- ❌ Basic input validation only

**Gaps:**
- [ ] CSRF tokens on state-changing operations
- [ ] Rate limiting on auth endpoints
- [ ] Security headers (HSTS, X-Frame-Options, CSP)
- [ ] XSS prevention (template escaping - check existing)
- [ ] SQL injection prevention (sqlc helps, verify)
- [ ] Session fixation prevention
- [ ] Brute-force protection on login

**Launch Blocker:** 🟡 **PARTIAL** - CSRF and rate limiting are high priority

**Estimated Effort:** **M (2-3 days / 16-24 hours)**

**Critical Security Fixes:**
```
1. CSRF Protection (High Priority)
   - Use gorilla/csrf middleware
   - Add tokens to all forms
   - Verify on POST/PUT/DELETE

2. Rate Limiting (High Priority)
   - Limit login attempts: 5 per 15 minutes
   - Limit password reset: 3 per hour
   - Limit API endpoints: 100 req/min per user

3. Security Headers (Medium Priority)
   - Strict-Transport-Security (HSTS)
   - X-Frame-Options: DENY
   - X-Content-Type-Options: nosniff
   - Content-Security-Policy (start basic)
```

**Current Vulnerabilities:**
```
Assessed Risks:
1. CSRF attacks on job/estimate operations - MEDIUM RISK
   Impact: Attacker could create/delete jobs
   Mitigation: Add CSRF middleware

2. Brute force password attacks - MEDIUM RISK
   Impact: Account takeover
   Mitigation: Rate limiting + account lockout

3. Session hijacking - LOW RISK (if using HTTPS)
   Impact: Account access
   Mitigation: Secure cookie flags + SameSite
```

**Good News:**
- Using sqlc reduces SQL injection risk
- html/template auto-escapes, reducing XSS risk
- Structured logging aids security monitoring

---

### 8. Database & Data Management

**Status:** 🟡 **PARTIAL** - SQLite works for MVP, PostgreSQL needed for production

**Current State:**
- ✅ Goose migrations working
- ✅ sqlc for type-safe queries
- ✅ Foreign key constraints
- ⚠️ SQLite database (not production-ready for SaaS)
- ❌ No backup strategy
- ❌ No data export for customers

**Gaps:**
- [ ] PostgreSQL migration
- [ ] Database backup automation
- [ ] Point-in-time recovery
- [ ] Customer data export (GDPR requirement)
- [ ] Database connection pooling
- [ ] Migration rollback strategy

**Launch Blocker:** 🟡 **PARTIAL** - PostgreSQL is strongly recommended

**Estimated Effort:** **M (2-3 days / 16-24 hours)**

**Why SQLite Won't Work for SaaS:**
```
Problems:
1. File-based locking (poor concurrency)
2. Single-server deployment only
3. Backup requires file copy (risky)
4. No point-in-time recovery
5. Write contention with multiple users
6. Hard to scale horizontally

When It Breaks:
- 10+ concurrent users
- Multiple server instances
- Heavy write operations (quote updates)
```

**PostgreSQL Migration Path:**
```
Good News:
- Using sqlc means minimal code changes
- Goose migrations work with PostgreSQL
- Most SQL is compatible

Changes Needed:
1. Update connection string
2. Change driver: "sqlite3" → "postgres"
3. Update migrations:
   - TEXT → VARCHAR or TEXT (both work)
   - REAL → NUMERIC(10,2)
   - datetime('now') → NOW()
4. Test all queries
5. Update Docker compose for Postgres

Effort: 2-3 days including testing
```

**Backup Strategy:**
```
Managed PostgreSQL (Recommended):
- Render.com: $7/mo for 1GB, auto-backups
- Supabase: $25/mo with backups + auth
- Digital Ocean: $15/mo managed Postgres

Self-Hosted:
- pg_dump nightly to S3
- Point-in-time recovery via WAL archiving
- Test restore process quarterly
```

---

### 9. Metrics & Monitoring

**Status:** 🟡 **PARTIAL** - Basic logging, no error tracking or metrics

**Current State:**
- ✅ Structured logging (slog)
- ✅ Request IDs
- ✅ Request duration logging
- ✅ Health check endpoint
- ❌ No error tracking (Sentry, etc.)
- ❌ No uptime monitoring
- ❌ No business metrics tracking

**Gaps:**
- [ ] Error tracking and alerting
- [ ] Uptime monitoring
- [ ] Database query performance
- [ ] Business metrics (signups, quotes created, MRR)
- [ ] Stripe webhook monitoring

**Launch Blocker:** ❌ **NO** - Can launch with basic logging

**Estimated Effort:** **S-M (1-2 days / 8-16 hours)**

**Alpha Launch Minimum:**
```
Must Have:
1. Sentry for error tracking ($0-26/mo)
2. Uptime monitor (UptimeRobot free tier)
3. Log aggregation (tail -f is fine for alpha)

Nice to Have:
4. Metrics dashboard (Grafana/Prometheus)
5. Business analytics (Mixpanel/Amplitude)
6. Query performance monitoring
```

**Quick Wins:**
```
1. Add Sentry (30 minutes)
   - Capture panics
   - Track errors
   - Alert on critical issues

2. Set up UptimeRobot (15 minutes)
   - Monitor /health endpoint
   - Email on downtime

3. Log aggregation (if multi-server)
   - Ship logs to CloudWatch/Papertrail
   - Searchable logs
```

**Business Metrics to Track:**
```
Critical:
- New signups per day
- Active subscriptions
- MRR (Monthly Recurring Revenue)
- Churn rate

Important:
- Quotes created per user
- Estimates sent per user
- Time to first quote (onboarding metric)
- Support requests per user

Don't Overcomplicate Early:
- Use Stripe dashboard for MRR
- Manual tracking in spreadsheet is fine
- Add proper analytics after 50+ users
```

---

### 10. Operational Readiness

**Status:** 🟡 **PARTIAL** - Docker setup exists, deployment process needs work

**Current State:**
- ✅ Multi-stage Dockerfile
- ✅ Caddy reverse proxy
- ✅ GitHub Actions CI/CD
- ✅ Health check endpoint
- 🟡 Manual deployment via Docker
- ❌ No secrets management
- ❌ No backup automation
- ❌ No graceful shutdown for in-progress operations

**Gaps:**
- [ ] Secrets management (env vars in production)
- [ ] Automated backups
- [ ] Blue-green or rolling deployments
- [ ] Database migration strategy in production
- [ ] Graceful shutdown handling
- [ ] Log rotation
- [ ] SSL/TLS certificate automation (Caddy does this, good!)

**Launch Blocker:** 🟡 **PARTIAL** - Secrets management is important

**Estimated Effort:** **M (2-3 days / 16-24 hours)**

**Deployment Checklist:**
```
Pre-Launch:
1. Environment-based secrets (not .env files in repo)
2. Database backup before each migration
3. Health check monitoring
4. SSL certificate auto-renewal (Caddy handles)
5. Log shipping to external service
6. Rollback plan documented

Deployment Process:
1. Run migrations on production DB
2. Build new Docker image
3. Push to registry
4. Pull on server
5. Stop old container
6. Start new container
7. Verify health check
8. Monitor error rates for 30 min
```

**Secrets Management:**
```
Development:
- .env files (current approach, fine)

Production Options:
1. Environment variables (simplest)
   - Set in Docker compose or systemd
   - Never commit to repo

2. Secrets management service
   - AWS Secrets Manager
   - HashiCorp Vault
   - Doppler (SaaS option)

Alpha Launch: Environment variables are fine
```

**Graceful Shutdown:**
```
Current: None
Needed:
- Catch SIGTERM signal
- Finish in-flight requests
- Close database connections
- Exit cleanly

Implementation (20 lines of code):
- Use http.Server.Shutdown(ctx)
- 30-second timeout
```

---

## Critical Architecture Changes Needed

### 1. From Single-User to Multi-Tenant

**Current Architecture:**
```
HTTP Request → Handler → Repository → SQLite
                ↓
            No user context
            No org scoping
            Global data access
```

**Required Architecture:**
```
HTTP Request → Auth Middleware → Tenant Middleware → Handler → Repository → PostgreSQL
                     ↓                    ↓
                 session.UserID      ctx.OrgID

Repository Layer:
- ALL queries include WHERE org_id = ?
- Tenant ID from context
- Panic if org_id missing (fail-safe)
```

**Implementation Pattern:**
```go
// Middleware to extract tenant
func TenantContext(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get user from session
        user := GetUserFromSession(r)

        // Add org_id to context
        ctx := context.WithValue(r.Context(), orgIDKey, user.OrgID)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Repository helper
func (q *Queries) orgID(ctx context.Context) string {
    orgID, ok := ctx.Value(orgIDKey).(string)
    if !ok {
        panic("org_id missing from context")
    }
    return orgID
}

// Every query becomes:
SELECT * FROM jobs WHERE org_id = $1 AND id = $2
```

**Migration Strategy for Existing Data:**
```sql
-- Add org_id columns
ALTER TABLE jobs ADD COLUMN org_id TEXT REFERENCES organizations(id);
ALTER TABLE clients ADD COLUMN org_id TEXT REFERENCES organizations(id);
-- ... repeat for all tables

-- For MVP data, create default org
INSERT INTO organizations (id, name) VALUES ('default-org', 'MVP Instance');
UPDATE jobs SET org_id = 'default-org';
UPDATE clients SET org_id = 'default-org';
-- ... repeat for all tables

-- Make org_id NOT NULL after backfill
ALTER TABLE jobs ALTER COLUMN org_id SET NOT NULL;
```

### 2. From SQLite to PostgreSQL

**Why This Matters:**
- SQLite file locking will cause problems with 5+ concurrent users
- Backups are hard (copy file while locked = corruption risk)
- Can't scale to multiple servers
- No row-level locking (entire database locks on writes)

**Migration Effort:** Medium (2-3 days)

**Code Changes:**
```diff
// main.go
-import _ "github.com/mattn/go-sqlite3"
+import _ "github.com/lib/pq"

-db, err := sql.Open("sqlite3", cfg.DatabasePath+"?_foreign_keys=on")
+db, err := sql.Open("postgres", cfg.DatabaseURL)

// Goose
-goose.SetDialect("sqlite3")
+goose.SetDialect("postgres")
```

**Migration File Changes:**
```diff
// Date functions
-created_at TEXT NOT NULL DEFAULT (datetime('now'))
+created_at TIMESTAMP NOT NULL DEFAULT NOW()

// Numeric precision
-surcharge_percent REAL
+surcharge_percent NUMERIC(10,2)

// Boolean type (SQLite uses INTEGER)
-active INTEGER NOT NULL DEFAULT 1
+active BOOLEAN NOT NULL DEFAULT true

// UUID generation (optional, can keep TEXT)
-id TEXT PRIMARY KEY
+id UUID PRIMARY KEY DEFAULT gen_random_uuid()
```

**Testing Strategy:**
1. Run PostgreSQL in Docker locally
2. Run all migrations
3. Run full test suite
4. Manually test all features
5. Load test with 10+ concurrent users

### 3. From Public to Authenticated Routes

**Current State:**
All routes are public. No protection.

**Required Changes:**
```go
// router/router.go
func Register(mux *http.ServeMux, h *Handler, authMux *Middleware) {
    // Public routes
    mux.HandleFunc("GET /health", h.Health)
    mux.HandleFunc("GET /", h.LandingPage)
    mux.HandleFunc("GET /sign/{token}", h.GetSignaturePage)
    mux.HandleFunc("POST /sign/{token}", h.SubmitSignature)

    // Auth routes
    mux.HandleFunc("GET /login", h.GetLogin)
    mux.HandleFunc("POST /login", h.PostLogin)
    mux.HandleFunc("POST /logout", h.PostLogout)
    mux.HandleFunc("GET /signup", h.GetSignup)
    mux.HandleFunc("POST /signup", h.PostSignup)
    mux.HandleFunc("GET /verify-email/{token}", h.VerifyEmail)

    // Protected routes - require auth + tenant context
    protected := authMux.Chain(
        RequireAuth,
        TenantContext,
        CSRF,
    )

    mux.Handle("GET /app/", protected(http.HandlerFunc(h.GetApp)))
    mux.Handle("GET /jobs", protected(http.HandlerFunc(h.ListJobs)))
    // ... all existing routes become protected
}
```

### 4. Settings and Company Profile: Per-Tenant

**Current State:**
- Settings table has single row (id='default')
- Company profile has single row (id='default')

**Problem:**
Every tenant needs their own settings and company profile.

**Solution:**
```sql
-- Change PK from 'default' to org_id
ALTER TABLE settings DROP CONSTRAINT settings_pkey;
ALTER TABLE settings ADD COLUMN org_id TEXT REFERENCES organizations(id);
UPDATE settings SET org_id = 'default-org' WHERE id = 'default';
ALTER TABLE settings ADD PRIMARY KEY (org_id);

-- Same for company_profile
ALTER TABLE company_profile DROP CONSTRAINT company_profile_pkey;
ALTER TABLE company_profile ADD COLUMN org_id TEXT REFERENCES organizations(id);
UPDATE company_profile SET org_id = 'default-org' WHERE id = 'default';
ALTER TABLE company_profile ADD PRIMARY KEY (org_id);
```

**Code Changes:**
```go
// Queries change from:
SELECT * FROM settings WHERE id = 'default'

// To:
SELECT * FROM settings WHERE org_id = $1
```

---

## Prioritized Roadmap

### Phase 0: Foundation (MUST COMPLETE FIRST)

**Duration:** 1-2 weeks
**Effort:** 40-60 hours
**Goal:** Establish multi-tenant architecture

| Task | Effort | Priority | Dependencies |
|------|--------|----------|--------------|
| PostgreSQL migration | M (2-3 days) | P0 | None |
| Add organizations table | S (4 hours) | P0 | PostgreSQL |
| Add users table with auth | M (2 days) | P0 | Organizations |
| Add org_id to all tables | M (2 days) | P0 | Organizations |
| Update all queries for tenant scoping | L (3-4 days) | P0 | org_id columns |
| Tenant context middleware | S (4 hours) | P0 | Users table |
| Migration for existing MVP data | S (4 hours) | P0 | org_id columns |

**Acceptance Criteria:**
- [ ] PostgreSQL running in production
- [ ] All tables have org_id
- [ ] All queries scoped to organization
- [ ] Middleware extracts org from session
- [ ] MVP data migrated to default org
- [ ] Tests pass with multi-tenant data

**Why This First:**
Everything else depends on multi-tenancy. Building auth without org context means rebuilding later.

---

### Phase 1: Alpha Launch Minimum

**Duration:** 4-6 weeks (after Phase 0)
**Effort:** 120-160 hours
**Goal:** Accept first paying customer

#### Week 1-2: Authentication & User Management

| Task | Effort | Priority |
|------|--------|----------|
| User registration flow | M (2 days) | P0 |
| Email verification | M (1-2 days) | P0 |
| Login/logout | S (1 day) | P0 |
| Password reset flow | M (2 days) | P0 |
| Session management | M (1-2 days) | P0 |
| Auth middleware | S (4 hours) | P0 |
| Protected route enforcement | S (4 hours) | P0 |

**Libraries to Add:**
- `golang.org/x/crypto/argon2` (password hashing)
- `gorilla/sessions` or custom JWT (session management)

#### Week 3-4: Billing & Payments

| Task | Effort | Priority |
|------|--------|----------|
| Stripe SDK integration | S (4 hours) | P0 |
| Customer creation on signup | M (1 day) | P0 |
| Subscription creation (single plan) | M (1 day) | P0 |
| Webhook endpoint + verification | M (2 days) | P0 |
| Subscription status tracking | M (1 day) | P0 |
| Access gating based on status | M (1 day) | P0 |
| Trial period logic | S (4 hours) | P0 |
| Canceled subscription handling | M (1 day) | P1 |

**Stripe Events to Handle:**
- `customer.subscription.created`
- `customer.subscription.updated`
- `customer.subscription.deleted`
- `invoice.payment_succeeded`
- `invoice.payment_failed`

#### Week 5: Email & Marketing Site

| Task | Effort | Priority |
|------|--------|----------|
| Postmark integration | S (4 hours) | P0 |
| Email verification template | S (2 hours) | P0 |
| Password reset template | S (2 hours) | P0 |
| Payment receipt template | S (2 hours) | P0 |
| Landing page | M (2 days) | P0 |
| Pricing page | S (4 hours) | P0 |
| Signup page with Stripe | M (1 day) | P0 |
| Terms of Service | S (4 hours) | P0 |
| Privacy Policy | S (4 hours) | P0 |

#### Week 6: Security & Polish

| Task | Effort | Priority |
|------|--------|----------|
| CSRF protection | M (1 day) | P0 |
| Rate limiting on auth | M (1 day) | P0 |
| Security headers | S (4 hours) | P0 |
| Sentry error tracking | S (2 hours) | P1 |
| Uptime monitoring | S (1 hour) | P1 |
| Database backup automation | M (1 day) | P0 |
| Secrets management (env vars) | S (2 hours) | P0 |
| Deployment runbook | S (2 hours) | P1 |

**Alpha Launch Acceptance Criteria:**
- [ ] User can sign up with email + password
- [ ] Email verification required
- [ ] User can create account + start trial
- [ ] Stripe subscription created on signup
- [ ] User can enter payment info
- [ ] Subscription activates after trial
- [ ] Failed payments block access
- [ ] User can cancel subscription
- [ ] Canceled users can reactivate
- [ ] All data scoped to organization
- [ ] Legal pages published
- [ ] Basic landing page live
- [ ] HTTPS enforced
- [ ] Database backups running
- [ ] Error tracking active

---

### Phase 2: Post-Alpha Improvements

**Duration:** 2-4 weeks
**Effort:** 60-100 hours
**Goal:** Reduce friction, improve retention

| Feature | Effort | Priority | Business Impact |
|---------|--------|----------|-----------------|
| Onboarding tour | M (2 days) | P0 | Reduce time to first quote |
| Email estimate delivery | M (2 days) | P0 | Core workflow completion |
| PDF export for estimates | L (3-4 days) | P0 | Professional appearance |
| Signature request emails | M (1 day) | P1 | Complete e-signature flow |
| Account settings page | M (1-2 days) | P1 | User empowerment |
| Team member invites | L (3-4 days) | P2 | Expansion revenue |
| Billing portal (Stripe) | M (2 days) | P1 | Self-service support |
| Quote duplication | M (1 day) | P1 | Efficiency gain |
| Basic analytics dashboard | M (2 days) | P2 | User engagement visibility |

**Recommendation:** Focus on onboarding first. A confused user churns immediately.

---

### Phase 3: Growth & Scale

**Duration:** Ongoing
**Goal:** Handle 100+ organizations, expand features

| Feature | Effort | Priority |
|---------|--------|----------|
| Multiple pricing tiers | L (3-4 days) | P1 |
| Upgrade/downgrade flows | L (3-4 days) | P1 |
| Advanced reporting | L (4-5 days) | P2 |
| QuickBooks integration | XL (2 weeks) | P2 |
| Mobile app (React Native) | XL (4-6 weeks) | P3 |
| Multi-user collaboration | XL (3-4 weeks) | P2 |
| API for integrations | L (1 week) | P3 |

**Defer Until:**
- 50+ paying customers
- Clear feature demand from users
- Stable MRR to fund development

---

## Data Model Changes Required

### New Tables Needed

#### 1. organizations
```sql
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(63) UNIQUE, -- optional: for subdomain routing

    -- Stripe
    stripe_customer_id VARCHAR(255) UNIQUE,
    stripe_subscription_id VARCHAR(255),

    -- Subscription
    plan VARCHAR(50) NOT NULL DEFAULT 'trial', -- trial, starter, pro, business
    status VARCHAR(50) NOT NULL DEFAULT 'active', -- active, past_due, canceled, suspended
    trial_ends_at TIMESTAMP,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_organizations_subdomain ON organizations(subdomain);
CREATE INDEX idx_organizations_stripe_customer ON organizations(stripe_customer_id);
```

#### 2. users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Identity
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,

    -- Profile
    first_name VARCHAR(100),
    last_name VARCHAR(100),

    -- Status
    email_verified BOOLEAN NOT NULL DEFAULT false,
    email_verification_token VARCHAR(255),
    email_verification_expires TIMESTAMP,

    -- Password reset
    password_reset_token VARCHAR(255),
    password_reset_expires TIMESTAMP,

    -- Role
    role VARCHAR(50) NOT NULL DEFAULT 'owner', -- owner, admin, member

    -- Timestamps
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_org ON users(org_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_verification_token ON users(email_verification_token);
CREATE INDEX idx_users_reset_token ON users(password_reset_token);
```

#### 3. sessions (if using database sessions)
```sql
CREATE TABLE sessions (
    id VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,

    -- Session data
    data JSONB,

    -- Expiry
    expires_at TIMESTAMP NOT NULL,

    -- Metadata
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user ON sessions(user_id);
CREATE INDEX idx_sessions_expires ON sessions(expires_at);
```

### Existing Tables: Add org_id

**Migration Strategy:**
```sql
-- Step 1: Add nullable org_id columns
ALTER TABLE jobs ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
ALTER TABLE clients ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
ALTER TABLE estimates ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
ALTER TABLE item_templates ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;
ALTER TABLE price_imports ADD COLUMN org_id UUID REFERENCES organizations(id) ON DELETE CASCADE;

-- Step 2: Backfill for existing data (create default org first)
INSERT INTO organizations (id, name, plan, status)
VALUES ('00000000-0000-0000-0000-000000000001', 'MVP Instance', 'pro', 'active');

UPDATE jobs SET org_id = '00000000-0000-0000-0000-000000000001';
UPDATE clients SET org_id = '00000000-0000-0000-0000-000000000001';
-- ... repeat for all tables

-- Step 3: Make org_id NOT NULL
ALTER TABLE jobs ALTER COLUMN org_id SET NOT NULL;
ALTER TABLE clients ALTER COLUMN org_id SET NOT NULL;
-- ... repeat

-- Step 4: Add indexes
CREATE INDEX idx_jobs_org ON jobs(org_id);
CREATE INDEX idx_clients_org ON clients(org_id);
-- ... repeat
```

### Settings & Company Profile: Refactor

**Current:**
```sql
-- Single row per table
CREATE TABLE settings (
    id TEXT PRIMARY KEY DEFAULT 'default',
    ...
);
```

**New:**
```sql
-- One row per organization
CREATE TABLE settings (
    org_id UUID PRIMARY KEY REFERENCES organizations(id) ON DELETE CASCADE,
    default_surcharge_mode VARCHAR(50) NOT NULL DEFAULT 'stacking',
    default_surcharge_percent NUMERIC(10,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE company_profiles (
    org_id UUID PRIMARY KEY REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(50),
    zip VARCHAR(20),
    logo_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

---

## Technical Recommendations

### 1. Use Established Libraries, Don't Roll Your Own

**Authentication:**
- ✅ DO: Use `golang.org/x/crypto/argon2` for password hashing
- ❌ DON'T: Implement your own crypto
- ✅ DO: Use `gorilla/sessions` or proven JWT library
- ❌ DON'T: Build custom session management

**CSRF Protection:**
- ✅ DO: Use `gorilla/csrf`
- ❌ DON'T: Implement token generation yourself

**Rate Limiting:**
- ✅ DO: Use `golang.org/x/time/rate` or `github.com/ulule/limiter`
- ❌ DON'T: Build in-memory rate limiting (won't scale)

### 2. Tenant Isolation Pattern

**Recommended Approach: Context Propagation**

```go
// Define context key
type contextKey string
const orgIDKey contextKey = "org_id"

// Middleware to extract from session
func TenantMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        session := GetSession(r)
        user := GetUser(session)

        ctx := context.WithValue(r.Context(), orgIDKey, user.OrgID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Helper in repository
func orgIDFromContext(ctx context.Context) (string, error) {
    orgID, ok := ctx.Value(orgIDKey).(string)
    if !ok || orgID == "" {
        return "", errors.New("org_id missing from context")
    }
    return orgID, nil
}

// Use in queries (sqlc)
-- name: ListJobsByOrg :many
SELECT * FROM jobs
WHERE org_id = $1
ORDER BY created_at DESC;
```

**Benefits:**
- Compiler-enforced org scoping
- Impossible to forget org_id filter
- Context travels through call stack
- Easy to test

### 3. Migration Strategy: Blue-Green for Zero Downtime

**Deployment Steps:**
```
1. Run new migrations on production DB
2. Deploy new code (backwards compatible)
3. Verify health checks pass
4. Monitor error rates for 30 minutes
5. If issues: rollback by redeploying previous version
```

**Making Migrations Backwards Compatible:**
```sql
-- BAD: Breaking change
ALTER TABLE jobs DROP COLUMN customer_name;

-- GOOD: Backwards compatible
ALTER TABLE jobs ADD COLUMN client_id UUID REFERENCES clients(id);
-- Deploy code that uses client_id
-- After 1 week, remove customer_name in separate migration
```

### 4. Stripe Integration Best Practices

**Webhook Security:**
```go
func HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
    payload, _ := io.ReadAll(r.Body)
    signature := r.Header.Get("Stripe-Signature")

    // VERIFY SIGNATURE (prevents fake webhooks)
    event, err := webhook.ConstructEvent(
        payload,
        signature,
        webhookSecret,
    )
    if err != nil {
        w.WriteHeader(400)
        return
    }

    // IDEMPOTENCY: Check if already processed
    if WebhookAlreadyProcessed(event.ID) {
        w.WriteHeader(200)
        return
    }

    // Process event
    switch event.Type {
    case "customer.subscription.updated":
        // Update subscription status in DB
    }

    // ACKNOWLEDGE immediately
    w.WriteHeader(200)
}
```

**Testing Webhooks:**
```bash
# Use Stripe CLI for local testing
stripe listen --forward-to localhost:8080/webhooks/stripe
stripe trigger customer.subscription.updated
```

### 5. Email Best Practices

**Deliverability:**
- Set up SPF, DKIM, DMARC records
- Use dedicated sending domain (mail.skalkaho.com)
- Monitor bounce rates
- Provide unsubscribe link (even for transactional)

**Template Management:**
```
Plain text is fine for alpha. HTML later.

Structure:
/email/templates/
  verify-email.txt
  password-reset.txt
  signature-request.txt
  payment-receipt.txt

Variables:
- {{.UserName}}
- {{.VerificationLink}}
- {{.EstimateURL}}
```

### 6. Security Headers (Middleware)

```go
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")

        next.ServeHTTP(w, r)
    })
}
```

---

## Risks & Mitigations

### Risk 1: Multi-Tenancy Bugs (Data Leaks)

**Risk:** Query without org_id exposes other customers' data

**Likelihood:** High (easy to forget one query)
**Impact:** Catastrophic (privacy violation, legal liability)

**Mitigations:**
1. **Code Review Checklist:** Every query must include org_id check
2. **Testing:** Write tests with multiple orgs, verify isolation
3. **Panic on Missing Context:** Helper function panics if org_id missing
4. **Database Constraints:** Consider row-level security (PostgreSQL RLS)
5. **Audit:** Regular SQL query audits for missing org_id

**Example Test:**
```go
func TestTenantIsolation(t *testing.T) {
    org1 := createOrg("Org 1")
    org2 := createOrg("Org 2")

    job1 := createJob(org1.ID, "Job 1")
    job2 := createJob(org2.ID, "Job 2")

    // Org 1 user should NOT see Org 2's jobs
    ctx := contextWithOrgID(org1.ID)
    jobs := ListJobs(ctx)

    assert.Len(t, jobs, 1)
    assert.Equal(t, job1.ID, jobs[0].ID)
}
```

### Risk 2: Payment Webhook Failures

**Risk:** Stripe webhook not received, subscription status out of sync

**Likelihood:** Medium (network issues, server downtime)
**Impact:** High (user blocked incorrectly, or unpaid user has access)

**Mitigations:**
1. **Idempotency:** Store webhook event IDs, ignore duplicates
2. **Retry Logic:** Stripe retries for 3 days
3. **Monitoring:** Alert on webhook failures
4. **Reconciliation:** Daily job to sync Stripe status with DB
5. **Grace Period:** Don't block immediately on failed payment

**Reconciliation Job:**
```go
// Run daily via cron
func ReconcileSubscriptions() {
    orgs := GetAllOrganizations()

    for _, org := range orgs {
        // Fetch from Stripe
        sub, _ := stripe.Subscription.Get(org.StripeSubscriptionID)

        // Compare with DB
        if sub.Status != org.Status {
            log.Warn("Status mismatch", "org", org.ID,
                "db_status", org.Status,
                "stripe_status", sub.Status)

            // Update DB to match Stripe (source of truth)
            UpdateOrgStatus(org.ID, sub.Status)
        }
    }
}
```

### Risk 3: Database Migration Failure in Production

**Risk:** Migration runs, fails halfway, database in inconsistent state

**Likelihood:** Low (if tested)
**Impact:** Critical (application down)

**Mitigations:**
1. **Test Migrations:** Run on copy of production data locally
2. **Backup First:** Automated backup before migration
3. **Transactions:** Wrap migrations in transactions (Goose supports)
4. **Rollback Plan:** Document rollback steps for each migration
5. **Maintenance Window:** Run during low-traffic periods

**Pre-Migration Checklist:**
```
[ ] Backup database
[ ] Test migration on production data copy
[ ] Document rollback steps
[ ] Verify migration is transactional
[ ] Schedule during low-traffic window
[ ] Have rollback SQL ready
[ ] Monitor error logs during migration
```

### Risk 4: Session Hijacking

**Risk:** Attacker steals session cookie, gains account access

**Likelihood:** Low (if HTTPS + secure cookies)
**Impact:** High (account takeover)

**Mitigations:**
1. **HTTPS Only:** Enforce TLS (Caddy does this)
2. **Secure Cookies:** HttpOnly, Secure, SameSite=Lax
3. **Session Expiry:** 30-day timeout
4. **IP Binding:** (Optional) Invalidate if IP changes drastically
5. **Activity Monitoring:** Alert on unusual activity

**Cookie Configuration:**
```go
store := sessions.NewCookieStore([]byte(cfg.SessionSecret))
store.Options = &sessions.Options{
    Path:     "/",
    MaxAge:   86400 * 30, // 30 days
    HttpOnly: true,        // Prevent JS access
    Secure:   true,        // HTTPS only
    SameSite: http.SameSiteLaxMode, // CSRF protection
}
```

### Risk 5: SQLite Concurrency Issues

**Risk:** File locking causes write conflicts with multiple users

**Likelihood:** High (if staying on SQLite)
**Impact:** Medium (errors, slow performance)

**Mitigations:**
1. **Migrate to PostgreSQL** (best solution)
2. If staying on SQLite temporarily:
   - WAL mode (Write-Ahead Logging)
   - Connection pool size = 1
   - Retry logic on SQLITE_BUSY
   - Monitor for lock timeout errors

**SQLite WAL Mode:**
```go
db, err := sql.Open("sqlite3", "file:quotes.db?cache=shared&mode=rwc&_journal_mode=WAL")
```

**Better Solution:** Migrate to PostgreSQL before launch.

---

## Alternative Approaches

### Option A: Full SaaS Launch (Recommended)

**Timeline:** 8-10 weeks
**Investment:** 200-250 hours
**Outcome:** Self-serve signups, Stripe billing, scalable

**Pros:**
- Can scale to 100+ customers immediately
- Automated billing reduces support burden
- Professional impression
- No manual provisioning

**Cons:**
- Longer time to first paying customer
- More upfront development
- Higher complexity

**Best For:**
- Confident in product-market fit
- Want to scale quickly
- Have development capacity
- Target self-serve market

### Option B: Pilot-First Launch (Faster Validation)

**Timeline:** 4-6 weeks
**Investment:** 100-150 hours
**Outcome:** 3-5 paying pilot customers, feedback loop

**Approach:**
1. Skip multi-tenancy (deploy separate instances per customer)
2. Add basic auth per instance (simple password)
3. Manual invoicing via Stripe invoices
4. Email delivery via Postmark
5. Personal onboarding calls

**Pros:**
- Launch in 4-6 weeks
- Learn from real customers faster
- Validate pricing
- Lower initial complexity
- Can pivot based on feedback

**Cons:**
- Doesn't scale (manual provisioning)
- Need to rebuild for true SaaS later
- Higher support burden
- Less professional

**Best For:**
- Uncertain about pricing
- Want customer feedback fast
- Small target market initially
- Willing to do manual work

**Transition Path:**
```
Weeks 1-6: Pilot launch (5 customers, manual)
Weeks 7-12: Build multi-tenant SaaS
Week 13: Migrate pilots to SaaS platform
Week 14+: Open self-serve signups
```

### Option C: Hybrid Approach

**Timeline:** 6-8 weeks
**Investment:** 150-200 hours

**Approach:**
1. Build multi-tenant architecture (Phase 0)
2. Build auth + basic billing
3. Manual onboarding initially
4. Automate onboarding later

**Pros:**
- Scalable architecture from start
- Can launch with manual onboarding
- Avoid rebuild later
- Lower risk

**Cons:**
- Still 6-8 weeks before first customer
- More work than Pilot-First
- Some features manual initially

**Best For:**
- Want scalability but okay with manual touches
- Have 6-8 weeks available
- Want to avoid technical debt

---

## Cost Breakdown (Alpha Launch)

### Development Costs

| Category | Hours | Rate ($100/hr) | Cost |
|----------|-------|----------------|------|
| Phase 0: Multi-tenancy | 40-60 | $100 | $4,000-6,000 |
| Phase 1: Auth + Billing | 80-100 | $100 | $8,000-10,000 |
| Phase 1: Email + Marketing | 30-40 | $100 | $3,000-4,000 |
| Phase 1: Security + Ops | 20-30 | $100 | $2,000-3,000 |
| **Total Development** | **170-230** | - | **$17,000-23,000** |

### Infrastructure Costs (Monthly)

| Service | Purpose | Cost |
|---------|---------|------|
| VPS (2 vCPU, 4GB RAM) | Application hosting | $12-24/mo |
| Managed PostgreSQL | Database | $7-15/mo |
| Postmark | Email delivery | $15/mo |
| Sentry | Error tracking | $0-26/mo (free tier works) |
| Domain + SSL | skalkaho.com | $12/year |
| Stripe | Payment processing | 2.9% + 30¢ per transaction |
| **Total Infrastructure** | - | **$35-70/mo** |

### Break-Even Analysis

**Assumptions:**
- Subscription price: $19/mo
- Development cost: $20,000 (amortized over 12 months = $1,667/mo)
- Infrastructure: $50/mo
- Monthly burn: $1,717/mo

**Break-Even:** 91 customers ($1,729 MRR)

**Cash-Flow Positive (ignoring dev cost):** 3 customers ($57 MRR > $50 costs)

---

## Success Metrics (First 90 Days)

### Leading Indicators (Measure Weekly)

| Metric | Week 4 Target | Week 8 Target | Week 12 Target |
|--------|---------------|---------------|----------------|
| Signups | 5 | 15 | 30 |
| Activated (created 1st quote) | 4 | 12 | 24 |
| Trial → Paid Conversion | 60% | 65% | 70% |
| Active Users (weekly) | 3 | 10 | 20 |
| Quotes Created per User | 2 | 3 | 4 |

### Lagging Indicators (Measure Monthly)

| Metric | Month 1 | Month 2 | Month 3 |
|--------|---------|---------|---------|
| MRR | $100 | $300 | $600 |
| Churn Rate | - | <10% | <5% |
| Support Tickets/User | <2 | <1 | <0.5 |
| NPS | 40+ | 50+ | 60+ |

### Red Flags (Immediate Action Required)

- **Activation Rate < 50%:** Onboarding is broken, users confused
- **Churn > 10%/month:** Product not delivering value
- **Support Tickets > 2/user/month:** Product too confusing or buggy
- **Time to First Quote > 1 hour:** Onboarding too complex
- **Trial → Paid < 50%:** Pricing wrong or value unclear

---

## Deployment Checklist

### Pre-Launch

**Infrastructure:**
- [ ] Production VPS provisioned
- [ ] PostgreSQL database created
- [ ] Database backups configured (daily)
- [ ] DNS configured (A record for skalkaho.com)
- [ ] SSL certificate provisioned (Caddy auto)
- [ ] Environment variables set
- [ ] Secrets management configured

**Application:**
- [ ] Multi-tenancy implemented and tested
- [ ] Authentication implemented and tested
- [ ] Stripe integration tested (use test mode first)
- [ ] Email delivery tested (Postmark)
- [ ] All migrations run successfully
- [ ] Seeds/fixtures for testing
- [ ] Health check responding

**Security:**
- [ ] CSRF protection enabled
- [ ] Rate limiting configured
- [ ] Security headers set
- [ ] HTTPS enforced (HTTP redirects)
- [ ] Session cookies secure
- [ ] Sentry error tracking active

**Legal:**
- [ ] Terms of Service published
- [ ] Privacy Policy published
- [ ] Cookie policy (if EU)
- [ ] GDPR compliance reviewed

**Monitoring:**
- [ ] Uptime monitoring configured
- [ ] Error tracking active (Sentry)
- [ ] Log aggregation working
- [ ] Database monitoring
- [ ] Stripe webhook monitoring

### Launch Day

- [ ] Switch Stripe to live mode
- [ ] Final database backup
- [ ] Deploy to production
- [ ] Smoke test all critical flows
- [ ] Announce to pilot list/network
- [ ] Monitor error rates for 24 hours

### Post-Launch (First Week)

- [ ] Daily check of error logs
- [ ] Monitor signup conversion
- [ ] Review first user feedback
- [ ] Fix critical bugs immediately
- [ ] Send welcome email to first customers
- [ ] Schedule onboarding calls

---

## Appendix A: Package Recommendations

### Core Dependencies (Add These)

```go
// go.mod additions needed

// Authentication
"golang.org/x/crypto/argon2"           // Password hashing
"github.com/gorilla/sessions"          // Session management
"github.com/gorilla/csrf"              // CSRF protection

// Database
"github.com/lib/pq"                    // PostgreSQL driver

// Payments
"github.com/stripe/stripe-go/v78"     // Stripe SDK

// Email
"github.com/sendgrid/sendgrid-go"     // OR
"github.com/mattbaird/gochimp"        // OR
// Just use HTTP client to Postmark API (simplest)

// Rate Limiting
"golang.org/x/time/rate"               // OR
"github.com/ulule/limiter/v3"         // Full-featured

// Monitoring
"github.com/getsentry/sentry-go"      // Error tracking

// Validation
"github.com/go-playground/validator/v10" // Struct validation
```

### Current Dependencies (Keep These)

```
✅ github.com/google/uuid              // ID generation
✅ github.com/joho/godotenv            // .env loading
✅ github.com/pressly/goose/v3         // Migrations
✅ github.com/anthropics/anthropic-sdk-go // AI matching (keep for price import)
✅ github.com/xuri/excelize/v2         // Excel parsing
```

---

## Appendix B: Suggested File Structure

```
skalkaho/
├── cmd/
│   └── server/
│       ├── main.go
│       └── migrations/          # Database migrations
├── internal/
│   ├── auth/                    # NEW: Authentication logic
│   │   ├── password.go
│   │   ├── session.go
│   │   └── middleware.go
│   ├── billing/                 # NEW: Stripe integration
│   │   ├── stripe.go
│   │   ├── webhooks.go
│   │   └── subscriptions.go
│   ├── email/                   # NEW: Email service
│   │   ├── postmark.go
│   │   └── templates/
│   ├── domain/                  # KEEP: Business logic
│   ├── handler/
│   │   ├── auth/               # NEW: Auth handlers
│   │   ├── billing/            # NEW: Billing handlers
│   │   ├── marketing/          # NEW: Public pages
│   │   └── keyboard/           # KEEP: App handlers
│   ├── middleware/             # EXPAND: Add auth, CSRF, rate limit
│   ├── repository/             # EXPAND: Add org scoping
│   └── templates/
│       ├── email/              # NEW: Email templates
│       ├── marketing/          # NEW: Landing, pricing, legal
│       └── keyboard/           # KEEP: App templates
├── web/                        # Static assets
├── planning/                   # Planning docs (this file)
└── development/                # Development docs
```

---

## Appendix C: Migration from SQLite to PostgreSQL

### Step-by-Step Guide

**1. Install PostgreSQL Locally**
```bash
# macOS
brew install postgresql@16
brew services start postgresql@16

# Linux (Ubuntu)
sudo apt install postgresql-16
sudo systemctl start postgresql

# Docker (easiest for testing)
docker run --name skalkaho-postgres \
  -e POSTGRES_PASSWORD=devpassword \
  -e POSTGRES_DB=skalkaho_dev \
  -p 5432:5432 \
  -d postgres:16
```

**2. Update Dependencies**
```bash
go get github.com/lib/pq
```

**3. Update Connection Code**
```go
// internal/database/database.go
import _ "github.com/lib/pq"

func Open(cfg *config.Config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        return nil, err
    }

    // Connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    return db, nil
}
```

**4. Update Goose Dialect**
```go
// cmd/server/main.go
-goose.SetDialect("sqlite3")
+goose.SetDialect("postgres")
```

**5. Convert Migrations**

For each migration file, update SQLite-specific syntax:

```sql
-- SQLite
created_at TEXT NOT NULL DEFAULT (datetime('now'))
surcharge_percent REAL
parent_id TEXT REFERENCES categories(id)

-- PostgreSQL
created_at TIMESTAMP NOT NULL DEFAULT NOW()
surcharge_percent NUMERIC(10,2)
parent_id UUID REFERENCES categories(id)
```

**6. Test Migrations**
```bash
# Create test database
createdb skalkaho_test

# Run migrations
DATABASE_URL=postgres://localhost/skalkaho_test go run cmd/server/main.go

# Verify schema
psql skalkaho_test -c "\dt"
```

**7. Update sqlc Config**
```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"  # Change from sqlite
    queries: "sqlc/queries"
    schema: "migrations"
    gen:
      go:
        package: "repository"
        out: "internal/repository"
```

**8. Regenerate sqlc Code**
```bash
make sqlc
```

**9. Test Application**
```bash
# Run all tests
go test ./...

# Manual testing
DATABASE_URL=postgres://localhost/skalkaho_dev go run cmd/server/main.go
```

**10. Production Migration**

```bash
# Create production database (managed service recommended)
# Render.com, Supabase, or Digital Ocean

# Set DATABASE_URL
export DATABASE_URL=postgres://user:pass@host:5432/dbname

# Run migrations
go run cmd/server/main.go
```

### Common Migration Issues

**Issue 1: Boolean vs Integer**
```sql
-- SQLite (uses INTEGER)
active INTEGER NOT NULL DEFAULT 1

-- PostgreSQL (use BOOLEAN)
active BOOLEAN NOT NULL DEFAULT true
```

**Issue 2: TEXT vs VARCHAR**
```sql
-- Both work in PostgreSQL, but VARCHAR is more idiomatic
name TEXT       # Works
name VARCHAR(255)  # Better
```

**Issue 3: AUTOINCREMENT vs SERIAL**
```sql
-- SQLite
id INTEGER PRIMARY KEY AUTOINCREMENT

-- PostgreSQL
id SERIAL PRIMARY KEY
-- OR better (for UUIDs)
id UUID PRIMARY KEY DEFAULT gen_random_uuid()
```

**Issue 4: Date Functions**
```sql
-- SQLite
datetime('now')
date('now', '+7 days')

-- PostgreSQL
NOW()
NOW() + INTERVAL '7 days'
CURRENT_DATE
```

---

## Conclusion

Skalkaho has a **strong product foundation** with sophisticated domain logic for construction quoting. The core value proposition is proven, and the existing features are well-implemented.

**The path to alpha launch is clear:**

1. **Foundation (2 weeks):** Multi-tenancy + PostgreSQL
2. **Authentication (1 week):** User accounts + sessions
3. **Billing (1 week):** Stripe integration
4. **Email (3 days):** Transactional emails
5. **Marketing (3 days):** Landing + legal pages
6. **Security (3 days):** CSRF + rate limiting
7. **Ops (2 days):** Backups + monitoring

**Total: 8-10 weeks to first paying customer**

**Key Risks to Manage:**
- Multi-tenancy bugs (data leaks)
- Webhook delivery failures
- Database migration issues
- Time estimation (always add buffer)

**Success Factors:**
- Start with Phase 0 (multi-tenancy foundation)
- Don't skip security hardening
- Test with multiple orgs before launch
- Have rollback plans for migrations
- Monitor everything from day one

**Alternative:** Consider Pilot-First approach (Option B) if you want to validate pricing and market fit faster with 4-6 weeks of work instead of 8-10.

The infrastructure work ahead is **well-defined and tractable**. No architectural surprises. The hard problems (surcharge inheritance, e-signatures) are already solved. What remains is plumbing—important plumbing, but plumbing nonetheless.

**You're in a strong position to launch.** The product works. The market is real. The execution plan is clear. Time to build the SaaS infrastructure and ship it.

---

**Document Version:** 1.0
**Date:** February 4, 2026
**Next Review:** After Phase 0 completion
**Contact:** feedback@skalkaho.com (future)

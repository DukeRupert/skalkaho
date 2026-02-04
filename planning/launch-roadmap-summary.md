# Skalkaho Alpha Launch Roadmap - Executive Summary

**Date:** February 4, 2026
**Status:** 50% Ready for Launch (Phase 0 complete)
**Time to Launch:** 6-8 weeks
**Estimated Effort:** 150-180 hours remaining

---

## Current State

**Strengths:**
- Feature-rich MVP with sophisticated quoting engine
- E-signature workflow complete
- Client management working
- Price import with AI matching
- Excellent domain logic and calculations

**Critical Gaps:**
- ~~No multi-tenancy (single-user app)~~ **DONE**
- No authentication or user accounts
- No billing/payment integration
- ~~SQLite database (won't scale)~~ **DONE** (PostgreSQL support added)
- No email delivery

---

## Top 5 Launch Blockers

| # | Blocker | Effort | Status |
|---|---------|--------|--------|
| 1 | ~~Multi-tenancy~~ | ~~1-2 weeks~~ | **DONE** |
| 2 | User Authentication | 4-5 days | Not started |
| 3 | Stripe Billing | 4-5 days | Not started |
| 4 | ~~PostgreSQL Migration~~ | ~~2-3 days~~ | **DONE** |
| 5 | Email Infrastructure | 2-3 days | Not started |

---

## Recommended Roadmap

### Phase 0: Foundation (2 weeks) - **COMPLETE**
**Goal:** Multi-tenant architecture

- [x] Migrate to PostgreSQL (Task 0.1)
- [x] Add organizations table (Task 0.2)
- [x] Add users table (Task 0.2)
- [x] Add org_id to all existing tables (Task 0.3)
- [x] Update all queries for tenant scoping (Task 0.3)
- [ ] Tenant context middleware (deferred to Phase 1 Auth)

### Phase 1: Alpha Launch (4-6 weeks)
**Goal:** First paying customer

**Weeks 1-2: Authentication**
- User registration + email verification
- Login/logout + password reset
- Session management
- Protected routes

**Weeks 3-4: Billing**
- Stripe customer creation
- Subscription management (single plan)
- Webhook handling
- Trial period logic
- Access gating

**Week 5: Email & Marketing**
- Postmark integration
- Email templates (verification, reset, receipts)
- Landing page with pricing
- Signup flow
- Terms of Service + Privacy Policy

**Week 6: Security & Launch**
- CSRF protection
- Rate limiting
- Security headers
- Database backups
- Error tracking (Sentry)
- Deploy and monitor

### Phase 2: Post-Launch (2-4 weeks)
**Goal:** Reduce friction, improve retention

- Onboarding tour
- Email estimate delivery
- PDF export for estimates
- Account settings
- Billing portal

---

## Alternative: Pilot-First Launch (4-6 weeks)

**Faster path with manual processes:**

- Skip multi-tenancy (separate instances per customer)
- Basic auth per instance
- Manual invoicing via Stripe
- Personal onboarding calls
- 3-5 pilot customers
- Learn and iterate

**Then rebuild for true SaaS in weeks 7-12**

---

## Critical Decisions

### 1. Full SaaS vs. Pilot-First?

**Full SaaS (8-10 weeks):**
- Can scale immediately
- Self-serve signups
- Professional appearance
- More upfront work

**Pilot-First (4-6 weeks):**
- Faster to market
- Learn from real customers
- Validate pricing
- Need to rebuild later

**Recommendation:** Full SaaS if confident in business model, Pilot-First if want faster validation.

### 2. SQLite vs. PostgreSQL?

**DONE - PostgreSQL support added.**

The codebase now supports both SQLite (for development) and PostgreSQL (for production).
Use `DATABASE_URL=postgres://...` environment variable to use PostgreSQL.

---

## Budget

### Development Costs
- ~~Phase 0: $4,000-6,000 (40-60 hours)~~ **COMPLETE**
- Phase 1: $13,000-17,000 (130-170 hours)
- **Remaining: $13,000-17,000** (or 130-170 hours internal)

### Infrastructure (Monthly)
- VPS: $12-24
- PostgreSQL: $7-15
- Postmark: $15
- Sentry: $0-26 (free tier)
- Stripe: 2.9% + 30¢ per transaction
- **Total: $35-70/month**

### Break-Even
- At $19/month: Need 91 customers (ignoring dev costs)
- Cash-flow positive: 3 customers ($57 MRR > $50 costs)

---

## Key Risks

1. **Multi-tenancy bugs** - Test thoroughly, data leaks are catastrophic
2. **Stripe webhook failures** - Implement retry + reconciliation
3. **Migration failures** - Backup before migrations, test on prod data copy
4. **Session hijacking** - HTTPS + secure cookies + SameSite
5. **Time estimation** - Add 25% buffer to all estimates

---

## Success Metrics (First 90 Days)

### Week 12 Targets
- 30 signups
- 24 activated (created first quote)
- 70% trial-to-paid conversion
- 20 active weekly users
- 4 quotes per user per month
- $600 MRR
- <5% churn

### Red Flags
- Activation < 50% (onboarding broken)
- Churn > 10% (product not valuable)
- Support tickets > 2/user/month (too confusing)
- Trial-to-paid < 50% (pricing wrong)

---

## Next Steps

1. ~~**Decision:** Full SaaS or Pilot-First?~~ **Decided: Full SaaS**
2. ~~**Start:** Phase 0 - PostgreSQL migration~~ **COMPLETE**
3. **Now:** Phase 1 - Authentication (sessions, login/logout, protected routes)
4. **Then:** Phase 1 - Stripe Billing integration
5. **Then:** Phase 1 - Email infrastructure (Postmark)
6. **Launch:** Ship when Phase 1 complete
7. **Learn:** Monitor metrics, iterate

---

## Files

- Full Analysis: `/workspaces/skalkaho/planning/saas-launch-readiness-2026-02-04.md`
- This Summary: `/workspaces/skalkaho/planning/launch-roadmap-summary.md`

---

**Bottom Line:** You have a strong product. The SaaS infrastructure is well-defined work. 8-10 weeks to launch, or 4-6 weeks if you choose the pilot-first path. The hardest problems are already solved. Time to build the plumbing and ship.

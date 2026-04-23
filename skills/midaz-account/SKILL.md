---
name: midaz-account
version: 0.7.0
description: "Authenticates, redeems invitations, completes initial onboarding, and manages the user's Midaz subscription via the midaz CLI. Use this skill whenever the user needs to sign in, sign out, check auth state, redeem an invite code, or manage their subscription — even if they don't mention Midaz by name."
when_to_use: "Trigger on phrases like 'sign in', 'log in', 'login', 'sign out', 'log out', 'am I logged in', 'who am I', 'whoami', 'redeem invite', 'invite code', 'invitation', 'subscription', 'upgrade my plan', 'cancel subscription', 'billing', 'my account', 'account status', 'set up Midaz', 'activate my account', 'credentials'."
metadata: {"requires":{"bins":["midaz"]}}
---

# Midaz Account

> Read [midaz-shared](../midaz-shared/SKILL.md) for response format, auth model, and side-effect gating.

Everything an agent needs to get a user from "never used Midaz" to "fully onboarded and subscribed". Covers auth, invitations, onboarding, Stripe, and PAT management.

## First-Time Flow

A fresh user goes through these steps in order. Run `midaz auth status` first if you're unsure where they are.

```
1. midaz auth login                         # sign in (browser)
2. midaz invite redeem <CODE> --yes         # unlock desk access
3. midaz onboard generate --mode guided \   # set radar + playbook
     --from-file guided.json --yes
4. midaz subscription start --yes           # start Stripe Checkout (opens browser)
5. midaz desk telegram connect              # (optional) wire up alerts
```

Use `midaz desk get` after each step to verify the state you expect.

## Authentication

```
midaz auth login              # Browser OAuth + exchange (default)
midaz auth login --paste      # Paste a PAT you created on the website
midaz auth login --token sk_… # Inline PAT (CI / headless)
midaz auth logout             # Clear stored credentials
midaz auth status             # Live whoami via GET /api/app/me
midaz auth whoami             # Offline: print cached credentials
```

PAT management:

```
midaz auth keys list                        # List active PATs
midaz auth keys create --label <name> --yes # Create a new PAT (shown once)
midaz auth keys revoke <id> --yes           # Revoke by id
```

Env overrides:
- `MIDAZ_TOKEN=sk_…` — inline credential for CI.
- `MIDAZ_NO_BROWSER=1` — never auto-open browsers.
- `MIDAZ_AUTH_PATH=/path/auth.json` — override credentials file location.

If any authenticated command exits with code 6, re-run `midaz auth login`.

## Invitations

Midaz uses single-use invitation codes to gate early access. A new user's desk has `has_invite_access: false` until a code is redeemed.

```
midaz invite redeem <CODE> --yes
```

Returns `{ success: true, message: "…" }`. After redemption, `midaz desk get` will show `has_invite_access: true`.

Common failures:
- `400 already_used` — the code was already redeemed.
- `404 not_found` — the code does not exist.

## Onboarding

A desk is "onboarded" when it has both a **radar** (watchlist, ≤12 items) and a **playbook** (trading rules, ≤20 000 chars). Two paths:

### Guided/freeform (server generates for you)

```
midaz onboard generate --mode guided --from-file guided.json --yes
midaz onboard generate --mode freeform --from-file free.json --yes
```

**Guided input** (`guided.json`):
```json
{
  "trading_style":  ["discretionary"],
  "time_horizon":   ["swing"],
  "market_scope":   ["equities", "macro"],
  "focus_areas":    ["Fed policy", "AI capex", "oil"],
  "notes":          "optional free-form context, ≤2000 chars"
}
```

**Freeform input** (`free.json`):
```json
{
  "freeform": "I trade macro swings in equities and rates, focused on Fed policy and energy. 20–4000 chars here."
}
```

Server LLM generates a radar + playbook and commits them atomically.

### Complete (you supply both)

```
midaz onboard complete --radar radar.md --playbook playbook.md --yes
```

Use when the user already handed you written content. Validates length and commits via `POST /api/desk/onboard`.

### Check state

```
midaz onboard status   # reads `onboarded` from GET /api/desk
```

## Subscription (Stripe)

Midaz Prime is $39/mo with a 14-day trial. Subscription status is on the desk:

```
midaz subscription status   # reads subscription object from GET /api/desk
```

Start a trial / subscribe:

```
midaz subscription start --yes
# → envelope contains a Stripe Checkout URL; auto-opens in browser.
```

After completing Checkout, `midaz subscription status` will report `status: "trialing"` (or `active` for returning customers who skip the trial).

Manage billing (update card, cancel, download invoices):

```
midaz subscription portal --yes
# → Stripe Customer Portal URL; auto-opens.
```

Allowed subscription states for subscription-gated endpoints (like `desk view`, `intel push`): `trialing`, `active`, `past_due`, or `dev_unlimited: true`. Exit code 7 from any command means the subscription is not in an allowed state.

## Agent Recipes

### "Help me sign in"

1. `midaz auth status` — if it returns 0, you're already in; show the email/desk.
2. If it exits 6: `midaz auth login` and tell the user to complete browser sign-in.
3. If on SSH / headless: run `midaz auth login --paste` and explain the user should create a PAT on the website first.

### "Set me up from scratch"

1. Ensure logged in.
2. `midaz desk get` — inspect `has_invite_access`, `onboarded`, `subscription.status`.
3. If `has_invite_access: false` → ask the user for their invitation code, then `midaz invite redeem <CODE> --yes`.
4. If `onboarded: false` → ask whether they want guided / freeform / manual; build the input file; `midaz onboard generate` or `onboard complete`.
5. If `subscription.status` is null or `canceled` → confirm, then `midaz subscription start --yes` and share the checkout URL.
6. Finish by surfacing the desk `view_url` so they can jump to the map.

### "Rotate my CLI credentials"

1. `midaz auth keys create --label "cli:<host>:<date>" --yes` → copy the new PAT.
2. `midaz auth login --token <new_pat>` to switch over (stores the new PAT).
3. `midaz auth keys list` → find the old key's id.
4. `midaz auth keys revoke <id> --yes`.

### "Cancel my subscription"

Don't run `subscription portal` without consent. Confirm with the user, then `midaz subscription portal --yes` and instruct them to cancel in the hosted portal. The CLI does not directly cancel.

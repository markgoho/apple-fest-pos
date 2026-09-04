---
status: accepted
---

# The System Admin and Leader PIN gates carry no session

Each PIN-gated route ([#43](https://github.com/markgoho/apple-fest-pos/issues/43)) always renders a PIN form first; a correct POST renders the page content in that same response, and nothing is stored anywhere — no cookie, no server-side session, no in-memory unlock state. A reload or a fresh visit always shows the blank form again, and there is no lockout on wrong attempts.

The obvious path is a session: unlock once, stay unlocked for a while. It was rejected because a session is state to build and keep correct for two humans who each visit briefly over one weekend, in a codebase that otherwise avoids server-side state beyond the single SQLite connection. Retyping a 4-digit PIN on each visit costs seconds; a session store costs a design surface (expiry, storage, invalidation) for a gain that doesn't matter at this scale — anyone who types the PIN correctly is standing at the device either way.

The same reasoning drops lockout: the threat model is a stranger at the booth on an air-gapped LAN, not a remote attacker who can run unlimited automated guesses, so a wrong-attempt counter buys little and risks locking out the real Leader.

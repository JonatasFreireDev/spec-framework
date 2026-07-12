# Events System Context

Status: `draft documentation fixture`

| Boundary | Proposed responsibility | Evidence limitation |
| --- | --- | --- |
| Attendee client | Presents QR state and requests token generation | No application code |
| Organizer client | Captures QR and displays validation result | No application code |
| Server authority | Generates and validates tokens and records attendance | Specification only |
| Data store | Owns event attendance and uniqueness constraints | No schema or migration |
| Analytics and audit | Records product events and security-relevant outcomes | No runtime evidence |

Trust boundaries and deployment topology remain unverified until a real adopter implementation exists.

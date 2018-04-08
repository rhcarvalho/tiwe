# Tiwe Game Protocol

## Determining gameplay order

- Some implicit initial order among players is assumed (*P1*, *P2*, *P3*, *P4*).
- Each player, following the implicit order:
  1. Chooses a random sequence *s* of 8 bytes.
  2. Computes the BLAKE2b-256 hash *h* of *s*.
  3. Shares *h* with all other players.
- After all players have shared their *h* value, each player, in implicit order,
  shares their *s* value. This guarantees that players cannot change their
  original sequence *s* without affecting the hash and being detected by other
  players.
- Every player:
  1. Computes *t* = BLAKE2b-256(XOR(*s1*, *s2*, *s3*, *s4*))
  2. *P1* ← uint64 big-endian of *t*[0:8]
  3. *P2* ← uint64 big-endian of *t*[8:16]
  4. *P3* ← uint64 big-endian of *t*[16:24]
  5. *P4* ← uint64 big-endian of *t*[24:32]
- The new defined gameplay order goes from the player with the lowest to the
  highest number.

Notes:

- The same algorithm applies to less players (minimum 2).
- For more than 4 players, the algorithm can be extended by:
  1. Compute *t2* = BLAKE2b-256(*t*)
  2. *P5* ← uint64 big-endian of *t2*[0:8]
  3. *P6* ← uint64 big-endian of *t2*[8:16]
  4. *P7* ← uint64 big-endian of *t2*[16:24]
  5. *P8* ← uint64 big-endian of *t2*[24:32]
  6. Compute *t3* = BLAKE2b-256(*t2*) ...
- In the very unlike case of collisions (2 or more players end up assigned with
  the same number), for simplicity, they should agree to retain their relative
  implicit order.

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
  players (proof of commitment).
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


## Shuffling tiles

- The first player generates a list of tiles, shuffle them and encrypt each
  using a commutative encryption scheme.
- The next players shuffle the list and encrypt with their own keys.
- After all players shuffled and encrypted the list, the first player decrypts
  each tile with her initial key, and then encrypt each tile with a different
  tile-specific key.
- The next players do the same: decrypt with their initial key and encrypt with
  their tile-specific keys.

### Commutative Encryption Scheme

Use SRA with specific choice of N (large prime).

### Tile representation

SRA exposes plain text information through quadratic residue.

Tile representation is done such that all tiles are a quadratic residue modulo
N.

```
For every tile; Jacobi(tile) == 1
```


---------------------------------------------------------------



- Every player upon program start listens for *Requests to Play*.
- One player starts a new game by sending a *Request to Play* to one or more peers.

- Upon receiving a *Request to Play*, players decide if they want to join the game and send a yes/no response.
- A game starts when all players reply to the *Request to Play* before the deadline.

- External parties may ask to *watch* a game by sending a *Request to Watch* to any player participating in a game.
- All active players and watchers maintain a Raft consensus (etcd) about the state of the game.
  Players form an ad-hoc etcd cluster, while watchers use the watch API to observe game state changes.


## Request to Play

```json
{
  "GameID": "xxxyyyzzz",
  "Deadline": "1985-04-12T23:20:50.52Z",
  "Players": ["alice", "bob", "carol"]
}
```

- Quorum: minimum number of players to start a game.
- Deadline: expiration timestamp.
  A player requesting to play will wait at most until the deadline to assume a
  negative answer from the players it didn't hear from.


## Starting the Game

- Initial order is based on the order in the *RequestToPlay*.
- Decide order based on **Determining gameplay order**.
- First player puts all tiles on the pool (encrypted) and shuffle.
- Next players reencrypt and shuffle.
- See https://en.wikipedia.org/wiki/Mental_poker#The_algorithm.

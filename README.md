# Deck

Deck of cards simulator

## API

All responses need to be JSON.

Described on [the OpenAPI spec](./api-spec.yaml).

## Design choices

- While I didn't want to set up a database for this projects, I did create a 
store API of sorts, that should be fairly straightforward to change to use a 
database.
- The tests are intentionally un-DRY, to make avoid any logic issues within the
tests themselves, and to make them easier to read.

## Compromises

- Didn't set up an option for creating a deck with a select set of cards.
Although all the required business logic is there to do so, it's missing that
parameter handling.
- Didn't populated decks with cards with `value` and `suit`, and relied on the
`code` attribute only.
- I wanted to setup a Swagger endpoint for the project but I ended up dropping
that idea due to time constraints.

## Issues

- The concurrent `handlePostDeckDraw` test is failing when ran with `-race`,
and I don't understand why at the moment. I've tried replicating the issue by
doing concurrent requests to the server but it never occurred. I wonder if the 
test itself is causing a race condition?


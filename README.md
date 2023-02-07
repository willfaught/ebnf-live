# ebnf-live

## Conflict endpoint

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' --basic -d '{"data":{"grammar":"..."}}' -u 'will:banana' $prefix/v1/conflict | jq
```

Example grammars:

- `S = A | B. A = z. B = z.`
- `S = A b. A = b | "".`

## First endpoint

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' -d '{"data":{"grammar":"..."}}' --basic -u 'will:banana' $host/v1/first | jq
```

Example grammars:

- `S = s | T. T = t.`
- `S = a A. A = b.`

## Follow endpoint

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' -d '{"data":{"grammar":"..."}}' --basic -u 'will:banana' $host/v1/follow | jq
```

Example grammars:

- `S = S a | S b | x.`
- `S = S s T. T = t.`

## Validate endpoint

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' -d '{"data":{"grammar":"..."}}' --basic -u 'will:banana' $host/v1/validate | jq
```

Example grammars:

- `S = s.`
- `s = s.`

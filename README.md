# ebnf-live

## Authorization

No authorization:

```sh
curl -X POST -i -s http://localhost/v1/validate
```

Bad authorization:

```sh
curl -X POST --basic -i -s -u 'will:bad' http://localhost/v1/validate
```

## Format

No content type:

```sh
curl -X POST --basic -i -s -u 'will:banana' http://localhost/v1/validate
```

Bad content type:

```sh
curl -X POST -H 'Content-Type: application/xml; charset=utf-16' --basic -i -s -u 'will:banana' http://localhost/v1/validate
```

## Endpoint

Bad path:

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' --basic -d '{"data":{"grammar":"S = s."}}' -i -s -u 'will:banana' http://localhost/v1
```

Bad method:

```sh
curl -X PUT -H 'Content-Type: application/json; charset=utf-8' --basic -d '{"data":{"grammar":"S = s."}}' -i -s -u 'will:banana' http://localhost/v1/validate
```

Bad JSON:

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' --basic -d 'bad' -i -s -u 'will:banana' http://localhost/v1/validate
```

Bad data:

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' --basic -d '[]' -i -s -u 'will:banana' http://localhost/v1/validate
```

## Conflict endpoint

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' --basic -d '{"data":{"grammar":"S = A | B. A = z. B = z."}}' -s -u 'will:banana' http://localhost/v1/conflict | jq
```

Example grammars:

- `S = A | B. A = z. B = z.`
- `S = A b. A = b | "".`

## First endpoint

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' --basic -d '{"data":{"grammar":"S = s | T. T = t."}}' -s -u 'will:banana' http://localhost/v1/first | jq
```

Example grammars:

- `S = s | T. T = t.`
- `S = a A. A = b.`

## Follow endpoint

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' --basic -d '{"data":{"grammar":"S = S a | S b | x."}}' -s -u 'will:banana' http://localhost/v1/follow | jq
```

Example grammars:

- `S = S a | S b | x.`
- `S = S s T. T = t.`

## Validate endpoint

```sh
curl -X POST -H 'Content-Type: application/json; charset=utf-8' --basic -d '{"data":{"grammar":"S = s."}}' -s -u 'will:banana' http://localhost/v1/validate | jq
```

Example grammars:

- `S = s.`
- `S = s. T = t.`

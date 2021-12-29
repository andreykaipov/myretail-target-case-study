# myRetail

This project contains two solutions to the Target _myRetail_ case study. The
prompt is copied over to [`PROMPT.md`](PROMPT.md) for convenience, but the TLDR
is to create a products API, that piggy-backs off of Target's own API.

## solution (traditional)

The first solution is your traditional 12-factor microservice -- written in Go,
using MongoDB as the NoSQL store to fulfill the prompt's requirements, and Redis
for anything we might want to cache. You can run it all locally using the
included [`docker-compose.yml`](./docker-compose.yml) file, or deploy just the
app onto your favorite container orchestrator platform, configuring it via the
environment to point to your remote MongoDB and Redis services. All it's missing
is a message queue before it turns into over-engineered hot garbage. 👻

Oh? Did somebody say demo? Totally.

```console
❯ docker-compose up -d
Creating network "traditional_cache" with the default driver
Creating network "traditional_database" with the default driver
Creating traditional_myretail-api_1 ... done
Creating traditional_redis_1        ... done
Creating traditional_mongo_1        ... done
```

Let's query one of the IDs provided to us:

```console
❯ curl -s http://localhost:3000/products/13860428 | jq .
{
  "id": 13860428,
  "name": "The Big Lebowski (Blu-ray)",
  "current_price": null
}
```

This proxied the request to RedSky and extracted only the interesting bits.
Unfortunately, however, the endpoint we've been tasked to use
(`/redsky_aggregations/v1/redsky/case_study_v1`) didn't include the pricing
info (hence the `null` price), so `myRetail` must now painstakingly create their
own competitive prices! Let's do that:

```console
❯ body='{"current_price":{"value":3.50,"currency_code":"USD"}}'
❯ curl -si http://localhost:3000/products/13860428 -XPUT -d "$body" -H 'Content-Type: application/json'
HTTP/1.1 401 Unauthorized
...
```

Oh right. I added authentication because _myRetail_ will never get to series
C funding if they have _yet another_ security scandal. However, if we check out
that [`docker-compose.yml`](./traditional/docker-compose.yml) file, we'll find
the expected creds. It was only a matter of time...

```console
❯ body='{"current_price":{"value":3.50,"currency_code":"USD"}}'
❯ curl -su admin:admin http://localhost:3000/products/13860428 -XPUT -d "$body" -H 'Content-Type: application/json'
HTTP/1.1 204 No Content
...
```

How's the GET looking now?

```console
❯ curl -s http://localhost:3000/products/13860428 | jq .
{
  "id": 13860428,
  "name": "The Big Lebowski (Blu-ray)",
  "current_price": {
    "value": 3.50,
    "currency_code": "USD"
  }
}
```

Nice. The final notable thing to showcase is how caching improves performance
here. What are we caching? Responses from RedSky! Not only are we responsible
consumers of its API, but the performance improvement speaks for itself:

```console
❯ curl -so/dev/null -w '%{time_total}s\n' http://localhost:3000/products/54456119
0.269485s
❯ curl -so/dev/null -w '%{time_total}s\n' http://localhost:3000/products/54456119
0.001669s
❯ curl -so/dev/null -w '%{time_total}s\n' http://localhost:3000/products/54456119
0.001504s
```

Subsequent requests are quicker by over 200ms! If we were to clear the cache,
our app will have to reach out to RedSky again.

```console
❯ redis-cli -n 0 hdel 54456119 name
(integer) 1
❯ curl -so/dev/null -w '%{time_total}s\n' http://localhost:3000/products/54456119
0.275151s
```

It's interesting to note all RedSky's APIs respond with the `cache-control`
header set to `max-age=0, no-cache, no-store`. So they definitely don't want
their CDN caching any page. However, since I only care about the product's name,
I feel pretty comfortable with caching it.

And that's about it for this approach. It's your typical cheese pizza. It's your
BLT. It's your chicken noodle soup.

## solution (serverless)

The second solution is a JavaScript service worker deployed onto Cloudflare's
serverless platform -- Cloudflare Workers. Demoing this is even nicer because
it's already running at https://myretail.kaipov.com! For deploying it onto your
own Cloudflare Worker though, see [`SETUP.md`](./serverless/traditional).

Let's prepopulate _myRetail_ with some test product data. We can do this by
hitting RedSky's `/redsky_aggregations/v1/web/plp_search_v1` endpoint. This was
found by just inspecting the network requests our browser makes when browsing
Target.com. The [`./fetch-tcins.sh`](./fetch-tcins.sh) script makes this pretty
easy. Just give it a category, and get some TCINs:

```console
❯ ./fetch-tcins.sh 5xtuv 3
82260756
82260792
82260695
```

How do we find categories though? Interestingly, it doesn't look like Target
exposes an API to list those. Instead available categories might be server-side
rendered and served directly in the HTML. And because I'm an **absolute
maniac**, I parse this raw HTML with even more shell:

```sh
curl -sL 'https://target.com/c/shop-all-categories/-/N-5xsxf' |
        sed 's!</a>!\n!g' |
        grep -Eo 'href="/c/.+?/-/N-[^"]+?"' |
        awk -F'[/?"]' '{print $6}' |
        cut -d- -f2 |
        xargs -I% ./fetch-tcins.sh % 25 |
        tee tcins.txt
```

This fetches the first 25 TCINs of each category, if the category even has TCINs
to begin with, and stores them in `tcins.txt`. Let's now populate our _myRetail_
KV store with randomly generated price info.

```sh
while read -r tcin; do
        price="$(seq 0 .01 20 | shuf | head -n1)"
        body="$(printf '{"current_price":{"value":%s,"currency_code":"USD"}}' "$price")"
        curl -su admin:password "https://myretail.kaipov.com/products/$tcin" \
                -XPUT -d "$body" \
                -H 'Content-Type: application/json'
done <tcins.txt
```

Awesome. I got through approximately half of that list before exceeding
Cloudflare's free plan limits! 😅 But... we can browse all of those tasty
products now:

```console
❯ curl -s https://myretail.kaipov.com/products?limit=5 | jq .
{
  "next_cursor": "AAAAAC-uOEjRSdsZu5ZEJ9Z6rJ4GXi2OMBkp_..."
  "ids": [
    10275315,
    10292372,
    10804811,
    10805587,
    10997227
  ]
}

❯ curl -s https://myretail.kaipov.com/products/10997227 | jq .
{
  "id": 10997227,
  "name": "Brita Replacement Water Filters for Brita Water Pitchers and Dispensers - 4ct",
  "current_price": {
    "value": 2.58,
    "currency_code": "USD"
  }
}
```

I can't believe it. $2.58!? What a great deal. This service also caches RedSky
responses. We can check it out ourselves when the cache is a hit or a miss by
inspecting the `redsky-cached` header:

```console
❯ curl -si https://myretail.kaipov.com/products/10805587 | grep redsky-cached
redsky-cached: MISS
❯ curl -si https://myretail.kaipov.com/products/10805587 | grep redsky-cached
redsky-cached: HIT
❯ curl -si https://myretail.kaipov.com/products/10805587 | grep redsky-cached
redsky-cached: HIT
```

Cool.

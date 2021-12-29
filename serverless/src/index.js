import {
    error,
    json,
    status,
    ThrowableRouter as Router,
    withContent,
    withParams,
} from 'https://esm.sh/itty-router-extras'

import { withBasicAuth } from './extras.js'
import * as redsky from './redsky.js'

const router = Router()
const withAuth = withBasicAuth(`${creds}`)

router.get('/products', async ({ query }) => {
    const data = await myretail.list({
        cursor: query.cursor,
        limit: query.limit || 50,
    })

    return json({
        next_cursor: data.cursor,
        ids: data.keys.map((x) => parseInt(x.name)),
    })
})

router.get('/products/:id', withParams, async ({ id }) => {
    const resp = await fetch(redsky.productURL(id), {
        cf: {
            cacheTtl: 60 * 60 * 24, // 24 hours
            cacheEverything: true,
        },
    })

    const parsed = await resp.json()
    const headers = {
        'redsky-cached': resp.headers.get('cf-cache-status'),
    }

    if (resp.status / 100 != 2) {
        return json(parsed, { status: resp.status, headers })
    }

    const priceInfo = await myretail.get(id, {})

    return json(
        {
            id: parseInt(id),
            name: parsed?.data?.product?.item?.product_description?.title,
            ...JSON.parse(priceInfo),
        },
        { headers },
    )
})

router.put('/products/:id', withParams, withContent, withAuth, async ({ id, content }) => {
    if (!content) {
        return error(400, 'Malformed or missing body')
    }

    if (content.id && content.id !== id) {
        return error(400, 'Mismatched IDs in path and PUT body')
    }

    await myretail.put(id, JSON.stringify(content))

    return status(204)
})

router.all('*', () => error(404, 'Lost in the sauce!'))

addEventListener('fetch', (event) => {
    event.respondWith(router.handle(event.request))
})

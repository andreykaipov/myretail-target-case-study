#!/bin/sh

usage() {
        cat <<EOF
Usage:
	$0 <category_id> [count] [offset]

Description:
	It's a tool to fetch TCINs from a Target product listing page.

Examples:
        $0 5xtuv 20 ....... Gets 20 products
        $0 5xtuv 5 10 ..... Gets 5 products, offset by the first 10
EOF
        exit 1
}

parseargs() {
        category="$1"
        count="${2-24}"
        offset="${3-0}"
        if [ -z "$category" ]; then usage; fi
}

main() {
        api=https://redsky-uat.perf.target.com/redsky_aggregations/v1/web
        key=3yUxt7WltYG7MFKPp7uyELi1K40ad2ys # this is not a secret key
        channel=WEB
        pricing_store_id=1083
        visitor_id=ghost

        curl -sG "$api/plp_search_v1" \
                -d "key=$key" \
                -d "channel=$channel" \
                -d "category=$category" \
                -d "page=/c/$category" \
                -d "pricing_store_id=$pricing_store_id" \
                -d "visitor_id=$visitor_id" \
                -d "count=$count" \
                -d "offset=$offset" | jq -r .data.search.products[].tcin
}

parseargs "$@"
main "$@"

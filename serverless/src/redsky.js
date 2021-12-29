const redsky = "https://redsky.target.com";
const productEndpoint = new URL(
  "/redsky_aggregations/v1/redsky/case_study_v1",
  redsky,
);
const key = "ff457966e64d5e877fdbad070f276d18ecec4a01";

export function productURL(id) {
  productEndpoint.search = new URLSearchParams({ key: key, tcin: id })
    .toString();
  return productEndpoint;
}

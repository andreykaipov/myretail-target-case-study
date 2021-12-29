import { error } from "https://esm.sh/itty-router-extras";

// Basic auth middleware for itty router

export const withBasicAuth = (creds) =>
  (request) => {
    const authHeader = request.headers.get("Authorization");

    if (!authHeader) {
      return error(401, "Credentials are required");
    }

    const expected = btoa(creds);
    const [scheme, credentials] = authHeader.split(" ");

    if (scheme !== "Basic" || credentials !== expected) {
      return error(401, "Credentials are incorrect");
    }
  };

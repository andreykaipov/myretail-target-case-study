with import <nixpkgs> { };

let
  name = "myretail";
in
stdenv.mkDerivation {
  name = "${name}-environment";
  buildInputs = [
    deno
    entr
    hey
    jq
    redis
    terraform
    terragrunt
    wrangler
  ];
  shellHook = ''
  '';
}

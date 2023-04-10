{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShell = pkgs.mkShell
          pkgs.mkShell
          {
            # The packages we need for this project
            buildInputs = with pkgs;
              [
                go_1_20
                go-tools
                gopls
              ];
          };
      });
}

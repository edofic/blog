{ nixpkgs ? import <nixpkgs> {}, compiler ? "default" }:

let

  inherit (nixpkgs) pkgs;

  f = { mkDerivation, base, hakyll, stdenv }:
      mkDerivation {
        pname = "blog";
        version = "0.1.0.0";
        src = ./generator;
        isLibrary = false;
        isExecutable = true;
        executableHaskellDepends = [ base hakyll ];
        license = stdenv.lib.licenses.asl20;
      };

  haskellPackages = if compiler == "default"
                       then pkgs.haskellPackages
                       else pkgs.haskell.packages.${compiler};

  drv = haskellPackages.callPackage f {};

in

  if pkgs.lib.inNixShell then drv.env else drv

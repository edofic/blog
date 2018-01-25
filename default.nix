{ stdenv
, glibcLocales
, haskell
}:

let

  haskellPackages = haskell.packages.ghc822;

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

in rec {

  blogBuilder = haskellPackages.callPackage f {}; # use .env for shell

  blogContent = stdenv.mkDerivation {
    name = "edofic-com";
    src = ./content;
    buildInputs = [ blogBuilder ];
    installPhase = ''
      export LANG="en_US.UTF-8"
      export LOCALE_ARCHIVE=${glibcLocales}/lib/locale/locale-archive

      blog build

      mkdir $out
      cp -R _site/. $out/.
    '';
  };

}

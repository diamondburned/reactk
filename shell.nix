let src = import ./nix/sources.nix {};
in

{ pkgs ? import src.nixpkgs {} }:

import "${src.gotk4-nix}/shell.nix" {
	pkgs = import "${src.gotk4-nix}/pkgs.nix" {
		sourceNixpkgs = src.nixpkgs;
		useFetched = true;
	};
	base = {
		pname = "gotk4-adwaita";
		version = "dev";

		buildInputs = pkgs: with pkgs; [
			gobject-introspection
			glib
			graphene
			gdk-pixbuf
			gtk4
			gtk3
			vulkan-headers
			libadwaita
		];
	};

	CGO_ENABLED = "1";
}

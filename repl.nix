{ pkgs }: {
    deps = [
        pkgs.go_1_20
        pkgs.chromium
        pkgs.gcc
        pkgs.pkg-config
    ];
}
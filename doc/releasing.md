# Publishing mcpcli Packages

This guide outlines the general steps to distribute `mcpcli` through common package managers. The GitHub release workflow builds cross-platform archives in the `dist/` directory and publishes packages to the configured targets.

## Homebrew

1. Create a tap repository such as `homebrew-mcpcli`.
2. Add a formula `mcpcli.rb` that references the macOS/Linux tarball from your GitHub release and its SHA256.
3. Commit and push the formula to the tap.
4. Users can install with:
   ```bash
   brew install <youruser>/mcpcli/mcpcli
   ```

## APT (Debian/Ubuntu)

1. Build a `.deb` package from the Linux binary using `dpkg-deb` or a tool like `fpm`.
2. Provide a `control` file describing the package name, version, maintainer and dependencies.
3. Sign the package and publish it via your own repository or a Launchpad PPA.
4. Users add the repository to their sources list and run `apt-get install mcpcli`.

## Chocolatey (Windows)

1. Create a nuspec file and install/uninstall PowerShell scripts.
2. Point the package to the Windows zip asset from the GitHub release.
3. Pack and test locally with `choco pack` and `choco install mcpcli`.
4. Push to the Chocolatey community repository with `choco push`.

The release workflow under `.github/workflows/release.yml` automatically updates the `homebrew`, `apt`, and `chocolatey` packages when a new tag is pushed.

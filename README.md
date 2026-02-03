# Sopsy 🔐

[![Release](https://img.shields.io/github/v/release/enbiyagoral/sopsy?style=flat-square)](https://github.com/enbiyagoral/sopsy/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/enbiyagoral/sopsy/ci.yaml?branch=main&style=flat-square)](https://github.com/enbiyagoral/sopsy/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/enbiyagoral/sopsy?style=flat-square)](https://goreportcard.com/report/github.com/enbiyagoral/sopsy)
[![License](https://img.shields.io/github/license/enbiyagoral/sopsy?style=flat-square)](https://github.com/enbiyagoral/sopsy/blob/main/LICENSE)


**Sopsy** is a lightweight profile manager for **SOPS**. Switch between encryption environments effortlessly, with automatic profile loading and sleek shell integration.

Currently supports **Age keys** and works with **Zsh** & **Bash**.

## Installation

```bash
brew install enbiyagoral/tap/sopsy
```

**Enable Shell Integration:**
_(Run once)_

```bash
# Zsh
sopsy init zsh
source ~/.zshrc

# Bash
sopsy init bash
source ~/.bashrc
```

---

## Usage

### Define Profiles

```bash
sopsy profile add stg --age-key-file ~/.config/sops/age/keys-stg.txt
sopsy profile add prod --age-key-file ~/.config/sops/age/prod.txt
```

### Switch Profiles

```bash
sopsy profile use stg
# ✓ Profile activated: stg
```

### Interactive Selection

If you have **fzf** installed, simply run:

```bash
sopsy profile use
```

> **Note:** Sopsy automatically loads your strict profile when you open a new terminal window.

---

## Use SOPS

```bash
sops -e -i secrets.yaml
sops -d secrets.yaml
```

## License

Apache-2.0.
#!/bin/bash

# Shared PATH for deploy scripts.
# Non-interactive shells (CI, bash script.sh) do not load ~/.bashrc,
# so pnpm installed via install.sh may be missing from PATH.

export PATH="$HOME/.local/share/pnpm:$HOME/.local/bin:/usr/local/bin:/usr/local/go/bin:$PATH"

if ! command -v pnpm >/dev/null 2>&1; then
	for candidate in \
		"$HOME/.local/share/pnpm/pnpm" \
		"/root/.local/share/pnpm/pnpm" \
		"/usr/local/bin/pnpm" \
		"/usr/bin/pnpm"; do
		if [ -x "$candidate" ]; then
			export PATH="$(dirname "$candidate"):$PATH"
			break
		fi
	done
fi

if ! command -v pnpm >/dev/null 2>&1 && command -v corepack >/dev/null 2>&1; then
	corepack enable pnpm >/dev/null 2>&1 || true
fi

if ! command -v pnpm >/dev/null 2>&1; then
	echo "pnpm не найден в PATH." >&2
	echo "Установите: curl -fsSL https://get.pnpm.io/install.sh | sh -" >&2
	echo "Или добавьте ~/.local/share/pnpm в PATH." >&2
	exit 1
fi

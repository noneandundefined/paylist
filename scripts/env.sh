#!/bin/bash

# Shared PATH for deploy scripts.
# Non-interactive shells (CI, bash script.sh) do not load ~/.bashrc,
# so pnpm from nvm or standalone install may be missing from PATH.

export PATH="$HOME/.local/share/pnpm:$HOME/.local/bin:/usr/local/bin:/usr/local/go/bin:$PATH"

_add_nvm_to_path() {
	local nvm_root="$1"

	if [ ! -d "$nvm_root/versions/node" ]; then
		return
	fi

	local node_bin
	for node_bin in "$nvm_root/versions/node"/*/bin; do
		if [ -d "$node_bin" ]; then
			export PATH="$node_bin:$PATH"
		fi
	done
}

_add_nvm_to_path "$HOME/.nvm"
_add_nvm_to_path "/root/.nvm"

if ! command -v pnpm >/dev/null 2>&1; then
	for candidate in \
		"$HOME/.local/share/pnpm/pnpm" \
		"/root/.local/share/pnpm/pnpm" \
		"/root/.nvm/versions/node/v24.19.0/bin/pnpm" \
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
	echo "Проверьте nvm: ls ~/.nvm/versions/node/*/bin/pnpm" >&2
	exit 1
fi

if ! command -v node >/dev/null 2>&1; then
	echo "node не найден в PATH (нужен для pnpm)." >&2
	exit 1
fi

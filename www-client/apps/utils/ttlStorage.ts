export const DAY_MS = 24 * 60 * 60 * 1000;

const memory = new Map<string, { expiresAt: number; value: unknown }>();

type TtlEntry<T> = {
	expiresAt: number;
	value: T;
};

const readEntry = <T>(key: string): TtlEntry<T> | null => {
	const fromMemory = memory.get(key);
	if (fromMemory) {
		return fromMemory as TtlEntry<T>;
	}

	if (typeof localStorage === 'undefined') {
		return null;
	}

	try {
		const raw = localStorage.getItem(key);
		if (!raw) {
			return null;
		}

		const parsed = JSON.parse(raw) as TtlEntry<T>;
		if (!parsed || typeof parsed.expiresAt !== 'number') {
			return null;
		}

		memory.set(key, parsed);
		return parsed;
	} catch {
		return null;
	}
};

export const readTtl = <T>(key: string, allowExpired = false): T | null => {
	const entry = readEntry<T>(key);
	if (!entry) {
		return null;
	}

	if (!allowExpired && Date.now() > entry.expiresAt) {
		return null;
	}

	return entry.value;
};

export const writeTtl = <T>(key: string, value: T, ttlMs: number): void => {
	const entry: TtlEntry<T> = {
		expiresAt: Date.now() + ttlMs,
		value,
	};

	memory.set(key, entry);

	if (typeof localStorage === 'undefined') {
		return;
	}

	try {
		localStorage.setItem(key, JSON.stringify(entry));
	} catch {
		return;
	}
};

export const withTtlCache = async <T>(key: string, ttlMs: number, loader: () => Promise<T>, isValid: (value: unknown) => value is T): Promise<T> => {
	const cached = readTtl<T>(key);
	if (cached !== null && isValid(cached)) {
		return cached;
	}

	try {
		const value = await loader();
		if (isValid(value)) {
			writeTtl(key, value, ttlMs);
		}

		return value;
	} catch (error) {
		const stale = readTtl<T>(key, true);
		if (stale !== null && isValid(stale)) {
			return stale;
		}

		throw error;
	}
};

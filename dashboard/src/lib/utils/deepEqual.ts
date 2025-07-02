// Deep equality check for objects/arrays
export function deepEqual(a: any, b: any): boolean {
	if (a === b) return true;
	if (typeof a !== typeof b) return false;
	if (typeof a !== 'object' || a === null || b === null) return false;
	if (Array.isArray(a) !== Array.isArray(b)) return false;

	if (Array.isArray(a)) {
		if (a.length !== b.length) return false;
		for (let i = 0; i < a.length; i++) {
			if (!deepEqual(a[i], b[i])) return false;
		}
		return true;
	}

	const keysA = Object.keys(a);
	const keysB = Object.keys(b);
	if (keysA.length !== keysB.length) return false;
	for (const key of keysA) {
		if (!deepEqual(a[key], b[key])) return false;
	}
	return true;
}

// Deep clone utility
export function deepClone<T>(obj: T): T {
	return JSON.parse(JSON.stringify(obj));
}

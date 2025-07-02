import { writable } from 'svelte/store';

export function persistentWritable<T>(key: string, initialValue: T, reviver?: (value: any) => T) {
    let storedValue: T | null = null;
    if (typeof localStorage !== 'undefined') {
        const json = localStorage.getItem(key);
        if (json) {
            const parsed = JSON.parse(json);
            storedValue = reviver ? reviver(parsed) : parsed;
        }
    }
    const store = writable<T>(storedValue ?? initialValue);

    store.subscribe((value) => {
        if (typeof localStorage !== 'undefined') {
            localStorage.setItem(key, JSON.stringify(value));
        }
    });

    return store;
} 
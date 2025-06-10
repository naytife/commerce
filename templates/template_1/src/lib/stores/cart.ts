import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface CartItem {
    id: string;
    title: string;
    price: number;
    quantity: number;
    image?: string;
    slug?: string;
}

function createCart() {
    const initial = browser ? JSON.parse(localStorage.getItem('cart') || '[]') as CartItem[] : [];
    const { subscribe, update, set } = writable<CartItem[]>(initial);

    if (browser) {
        subscribe(items => {
            localStorage.setItem('cart', JSON.stringify(items));
        });
    }

    return {
        subscribe,
        add: (item: Omit<CartItem, 'quantity'>, quantity = 1) => update(items => {
            const existing = items.find(i => i.id === item.id);
            if (existing) {
                return items.map(i => i.id === item.id ? { ...i, quantity: i.quantity + quantity } : i);
            }
            return [...items, { ...item, quantity }];
        }),
        remove: (id: string) => update(items => items.filter(i => i.id !== id)),
        updateQuantity: (id: string, quantity: number) => update(items => items.map(i => i.id === id ? { ...i, quantity } : i)),
        clear: () => set([])
    };
}

export const cart = createCart(); 
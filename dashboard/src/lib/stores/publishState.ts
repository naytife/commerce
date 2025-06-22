import { writable } from 'svelte/store';

export interface ChangeRecord {
  id: string;
  type: 'shop_update' | 'product_create' | 'product_update' | 'product_delete' | 'image_update';
  entity: string;
  timestamp: Date;
  description: string;
}

export interface PublishState {
  hasUnpublishedChanges: boolean;
  lastPublished: Date | null;
  changesSince: ChangeRecord[];
  isPublishing: boolean;
  publishError: string | null;
  currentVersion: string;
}

const initialState: PublishState = {
  hasUnpublishedChanges: false,
  lastPublished: null,
  changesSince: [],
  isPublishing: false,
  publishError: null,
  currentVersion: '1.0.0',
};

function createPublishState() {
  const { subscribe, set, update } = writable<PublishState>(initialState);

  return {
    subscribe,
    set,
    update,
    
    // Track a new change
    recordChange: (change: Omit<ChangeRecord, 'id' | 'timestamp'>) => {
      update(state => ({
        ...state,
        hasUnpublishedChanges: true,
        changesSince: [
          ...state.changesSince,
          {
            ...change,
            id: `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
            timestamp: new Date(),
          }
        ]
      }));
    },
    
    // Mark as published
    markAsPublished: () => {
      update(state => ({
        ...state,
        hasUnpublishedChanges: false,
        lastPublished: new Date(),
        changesSince: [],
        isPublishing: false,
        publishError: null,
      }));
    },
    
    // Set publishing state
    setPublishing: (isPublishing: boolean, error?: string) => {
      update(state => ({
        ...state,
        isPublishing,
        publishError: error || null,
      }));
    },
    
    // Clear all changes
    clearChanges: () => {
      update(state => ({
        ...state,
        hasUnpublishedChanges: false,
        changesSince: [],
      }));
    },
    
    // Get summary of changes
    getChangesSummary: () => {
      let summary = '';
      subscribe(state => {
        const counts = state.changesSince.reduce((acc, change) => {
          acc[change.type] = (acc[change.type] || 0) + 1;
          return acc;
        }, {} as Record<string, number>);
        
        const parts = [];
        if (counts.shop_update) parts.push(`${counts.shop_update} store update${counts.shop_update > 1 ? 's' : ''}`);
        if (counts.product_create) parts.push(`${counts.product_create} product${counts.product_create > 1 ? 's' : ''} added`);
        if (counts.product_update) parts.push(`${counts.product_update} product${counts.product_update > 1 ? 's' : ''} updated`);
        if (counts.product_delete) parts.push(`${counts.product_delete} product${counts.product_delete > 1 ? 's' : ''} deleted`);
        if (counts.image_update) parts.push(`${counts.image_update} image${counts.image_update > 1 ? 's' : ''} updated`);
        
        summary = parts.join(', ');
      })();
      return summary;
    }
  };
}

export const publishState = createPublishState();

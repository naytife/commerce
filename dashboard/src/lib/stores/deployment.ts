import { writable } from 'svelte/store';
import type { Shop } from '$lib/types';

export interface DeploymentStatus {
	shopId: string;
	subdomain: string;
	status: 'deploying' | 'deployed' | 'failed';
	message?: string;
	startedAt: Date;
	completedAt?: Date;
	url?: string; // Add URL field
}

export interface DeploymentStore {
	isDeploying: boolean;
	deployments: DeploymentStatus[];
}

const initialState: DeploymentStore = {
	isDeploying: false,
	deployments: []
};

function createDeploymentStore() {
	const { subscribe, set, update } = writable<DeploymentStore>(initialState);

	let pollInterval: NodeJS.Timeout | null = null;
	let authFetchFunction: typeof fetch | null = null;
	let pollAttempts = 0; // Track number of polling attempts for backoff

	return {
		subscribe,
		
		// Set the authenticated fetch function
		setAuthFetch: (authFetch: typeof fetch) => {
			authFetchFunction = authFetch;
		},
		
		// Start tracking a new deployment
		startDeployment: (shop: Shop) => {
			update(state => {
				const deployment: DeploymentStatus = {
					shopId: shop.shop_id.toString(),
					subdomain: shop.subdomain,
					status: 'deploying',
					startedAt: new Date()
				};
				
				return {
					isDeploying: true,
					deployments: [...state.deployments.filter(d => d.shopId !== shop.shop_id.toString()), deployment]
				};
			});
			
			// Reset poll attempts and start polling for this deployment
			pollAttempts = 0;
			startPolling(shop.shop_id.toString(), shop.subdomain);
		},
		
		// Stop tracking a deployment
		completeDeployment: (shopId: string, url?: string) => {
			update(state => {
				const updatedDeployments = state.deployments.map(d => 
					d.shopId === shopId 
						? { ...d, status: 'deployed' as const, completedAt: new Date(), url }
						: d
				);
				
				const stillDeploying = updatedDeployments.some(d => d.status === 'deploying');
				
				return {
					isDeploying: stillDeploying,
					deployments: updatedDeployments
				};
			});
			
			// Stop polling if no more deployments
			if (pollInterval) {
				clearTimeout(pollInterval);
				pollInterval = null;
				pollAttempts = 0;
			}
		},
		
		// Mark deployment as failed
		failDeployment: (shopId: string, message?: string) => {
			update(state => {
				const updatedDeployments = state.deployments.map(d => 
					d.shopId === shopId 
						? { ...d, status: 'failed' as const, message, completedAt: new Date() }
						: d
				);
				
				const stillDeploying = updatedDeployments.some(d => d.status === 'deploying');
				
				return {
					isDeploying: stillDeploying,
					deployments: updatedDeployments
				};
			});
			
			// Stop polling if no more deployments
			if (pollInterval) {
				clearTimeout(pollInterval);
				pollInterval = null;
				pollAttempts = 0;
			}
		},
		
		// Remove a deployment from tracking
		removeDeployment: (shopId: string) => {
			update(state => ({
				...state,
				deployments: state.deployments.filter(d => d.shopId !== shopId)
			}));
		},
		
		// Clear all deployments
		clearAll: () => {
			if (pollInterval) {
				clearTimeout(pollInterval);
				pollInterval = null;
				pollAttempts = 0;
			}
			set(initialState);
		},
		
		// Manually check deployment status (for user-triggered refresh)
		checkDeploymentStatus: async (shopId: string) => {
			if (!authFetchFunction) {
				console.warn('Auth fetch function not set');
				return;
			}
			
			try {
				const response = await authFetchFunction(`http://127.0.0.1:8080/v1/shops/${shopId}/deployment-status`);
				
				if (response.ok) {
					const data = await response.json();
					const deploymentData = data.data; // Extract the nested data
					
					if (deploymentData?.status === 'deployed') {
						deploymentStore.completeDeployment(shopId, deploymentData.url);
					} else if (deploymentData?.status === 'failed') {
						deploymentStore.failDeployment(shopId, deploymentData.message);
					}
				}
			} catch (error) {
				console.error('Error checking deployment status:', error);
			}
		}
	};

	async function startPolling(shopId: string, subdomain: string) {
		if (pollInterval) return; // Already polling
		
		const checkStatus = async () => {
			try {
				if (!authFetchFunction) {
					console.warn('Auth fetch function not set, skipping deployment status check');
					return;
				}

				const response = await authFetchFunction(`http://127.0.0.1:8080/v1/shops/${shopId}/deployment-status`);
				
				if (response.ok) {
					const data = await response.json();
					const deploymentData = data.data; // Extract the nested data
					
					if (deploymentData?.status === 'deployed') {
						deploymentStore.completeDeployment(shopId, deploymentData.url);
						return true; // Stop polling
					} else if (deploymentData?.status === 'failed') {
						deploymentStore.failDeployment(shopId, deploymentData.message);
						return true; // Stop polling
					}
				}
				return false; // Continue polling
			} catch (error) {
				console.error('Error checking deployment status:', error);
				return false; // Continue polling despite errors
			}
		};
		
		const scheduleNextPoll = () => {
			pollAttempts++;
			
			// Progressive intervals: start fast, then slow down
			// 1s, 1s, 2s, 3s, 5s, 5s, 5s...
			let interval = 1000; // Start with 1 second
			
			if (pollAttempts <= 2) {
				interval = 1000; // First 2 attempts: 1 second
			} else if (pollAttempts <= 4) {
				interval = 2000; // Next 2 attempts: 2 seconds  
			} else if (pollAttempts <= 6) {
				interval = 3000; // Next 2 attempts: 3 seconds
			} else {
				interval = 5000; // After that: 5 seconds
			}
			
			pollInterval = setTimeout(async () => {
				const shouldStop = await checkStatus();
				if (!shouldStop) {
					scheduleNextPoll();
				}
			}, interval);
		};
		
		// Start with immediate check, then schedule subsequent polls
		const shouldStop = await checkStatus();
		if (!shouldStop) {
			scheduleNextPoll();
		}
	}
}

export const deploymentStore = createDeploymentStore();

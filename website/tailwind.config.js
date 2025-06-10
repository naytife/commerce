/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			animation: {
				'custom-fade-1': 'custom-fade 0.6s  backwards',
				'custom-fade-2': 'custom-fade 2s  backwards',
				'custom-fade-3': 'custom-fade 1s  backwards',
				'custom-fade-4': 'custom-fade 2s  backwards',
				'custom-fade-5': 'custom-fade 1s  backwards',
				pop: 'pop 1s   forwards',
				// Existing animations
				'custom-fade': 'custom-fade 4s alternate linear',
				moveGradientright: 'slider 6s infinite linear',
				moveGradientleft: 'slide 6s infinite linear',
				// spinf and spins removed as unused
				// Add custom animations from your CSS
				'problem-item-box': 'problem-item-box 7s infinite backwards',
				'problem-item-label': 'problem-item-label 7s infinite backwards',
				'problem-item-icon': 'problem-item-icon 7s infinite backwards',
				'problem-item-icon-img': 'problem-item-icon-img 7s infinite backwards',
				'savings-logos-border': 'savings-optimization-logos-border 1s linear infinite',
				'fade-in-up-hero-headline': 'fade-in-up 0.8s ease-out forwards',
				'fade-in-up-hero-subheadline': 'fade-in-up 0.8s ease-out 0.3s forwards'
			},
			boxShadow: {
				'2xl': '0 10px 12px rgba(0, 0, 0, 0.1)',
				'problem-item': '0 0 0 1px #0e3f7e0f, 0 1px 1px -0.5px #2a334608, 0 2px 2px -1px #2a33460a, 0 3px 3px -1.5px #2a33460a, 0 5px 5px -2.5px #2a334608, 0 10px 10px -5px #2a334608, 0 24px 24px -8px #2a334608'
			},
			keyframes: {
				slide: {
					'0%': { backgroundPosition: '0% 0%' },
					'100%': { backgroundPosition: '-100% -100%' }
				},
				slider: {
					'0%': { backgroundPosition: '0% 0%' },
					'100%': { backgroundPosition: '100% -100%' }
				},
				// fadeIn keyframes removed as unused
				// Add keyframes from your CSS
				'problem-item-box': {
					'17.8%, 57.1%': { backgroundColor: 'var(--white-24)' },
					'25%, 50%': { backgroundColor: 'var(--white-72)' }
				},
				'problem-item-label': {
					'0%, 17.8%': { opacity: '0', transform: 'translateY(10px)' },
					'25%, 78.5%': { opacity: '1', transform: 'translateY(0)', filter: 'blur(0)' },
					'85.6%': {
						filter: 'blur(10px) grayscale(1)',
						opacity: '0',
						transform: 'translateY(-20px)'
					}
				},
				'problem-item-icon': {
					'0%, 14.2%': {
						opacity: '0',
						transform: 'translateY(20px)',
						filter: 'blur(10px) grayscale(1)'
					},
					'21.4%, 78.5%': {
						opacity: '1',
						transform: 'translateY(0)',
						filter: 'blur(0)'
					},
					'82%': {
						filter: 'blur(5px) grayscale(1)'
					},
					to: {
						opacity: '0',
						transform: 'translateY(-100px) scale(2)',
						filter: 'blur(20px) grayscale(1)'
					}
				},
				'problem-item-icon-img': {
					'0%, 17.8%, 83%': {
						height: '0px',
						width: '0px',
						transform: 'translateY(0)'
					},
					'32.1%, 78.5%': {
						height: '32px',
						width: '32px',
						opacity: '1',
						transform: 'translateY(0)'
					},
					'82%': {
						opacity: '0',
						height: '32px',
						width: '32px',
						transform: 'translateY(-10px)'
					}
				},
				pop: {
					'0%': { transform: 'scale(0.5) translateX(20%) translateY(20%)', opacity: '0' },
					'50%': { transform: 'scale(1.1) translateY(-10%)', opacity: '1' },
					'100%': { transform: 'scale(1) translateX(50%) translateY(-50%)', opacity: '1' }
				},
				'custom-fade': {
					'0%, 17.8%, 83%': {
						height: '0px',
						width: '0px',
						transform: 'translateY(0)',
						filter: 'blur(40px) grayscale(1)'
					},
					'32.1%, 78.5%': {
						height: '32px',
						width: '32px',
						opacity: '1',
						transform: 'translateY(0)',
						filter: 'blur(20px) grayscale(1)'
					},
					'82%': {
						opacity: '0',
						height: '50px',
						width: '50px',
						transform: 'translateY(-40px)',
						filter: 'blur(20px) grayscale(1)'
					}
				},
				'savings-optimization-logos-border': {
					'0%': { transform: 'translate(-50%, -50%) rotate(0deg)' },
					to: { transform: 'translate(-50%, -50%) rotate(360deg)' }
				},
				'fade-in-up': {
					'0%': { opacity: '0', transform: 'translateY(1rem)' },
					'100%': { opacity: '1', transform: 'translateY(0)' }
				}
			},
			colors: {
				// custom-blue removed as unused
				// Add more custom colors if needed
			},
			borderColor: {
				// neutral-opaque-6 removed as unused
			},
			// boxShadow is extended above
			borderRadius: {
				// xl: '1rem' removed as redundant
			},
			fontSize: {
				// sm: '0.875rem' removed as redundant
			},
			width: {
				// 150: '150px' removed as unused
				40: '10rem',
				168: '42rem' // Assuming w-168 refers to 42rem (optional, adjust as needed)
			},
			height: {
				// 150: '150px' removed as unused
				40: '10rem'
				// Add more custom heights if needed
			}
		}
	},
	plugins: []
};

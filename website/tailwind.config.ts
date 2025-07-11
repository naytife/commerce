import defaultTheme from 'tailwindcss/defaultTheme';
import type { Config } from 'tailwindcss';

const config: Config = {
	darkMode: 'class',
	content: ['./src/**/*.{html,js,svelte,ts}'],
	// safelist: ["dark"],
	theme: {
		container: {
			center: true,
			padding: '2rem',
			screens: {
				'2xl': '1400px'
			}
		},
		extend: {
			colors: {
				// Core brand colors - now using OKLCH format
				brand: {
					primary: 'oklch(var(--brand-primary))',
					secondary: 'oklch(var(--brand-secondary))',
					accent: 'oklch(var(--brand-accent))',
					// Additional brand variants in OKLCH
					50: 'oklch(0.95 0.05 141)',   // Very light green
					100: 'oklch(0.85 0.12 141)',  // Light green
					200: 'oklch(0.75 0.15 141)',  // Lighter green
					300: 'oklch(0.65 0.16 141)',  // Medium light green
					400: 'oklch(0.55 0.17 141)',  // Medium green
					500: 'oklch(var(--brand-primary))', // Primary green
					600: 'oklch(0.45 0.18 141)',  // Darker green
					700: 'oklch(0.35 0.16 141)',  // Dark green
					800: 'oklch(0.25 0.14 141)',  // Very dark green
					900: 'oklch(0.15 0.12 141)',  // Darkest green
					// Charcoal variants
					charcoal: {
						50: 'oklch(0.95 0.01 60)',   // Very light gray
						100: 'oklch(0.85 0.01 60)',  // Light gray
						200: 'oklch(0.75 0.01 60)',  // Medium light gray
						300: 'oklch(0.65 0.01 60)',  // Medium gray
						400: 'oklch(0.45 0.01 60)',  // Dark gray
						500: 'oklch(0.25 0.01 60)',  // Darker gray
						600: 'oklch(0.18 0.01 60)',  // Very dark gray
						700: 'oklch(0.15 0.01 60)',  // Charcoal
						800: 'oklch(var(--brand-secondary))', // Logo charcoal
						900: 'oklch(0.08 0.01 60)',   // Darkest
					}
				},
				// Base colors using OKLCH
				border: 'oklch(var(--border) / <alpha-value>)',
				input: 'oklch(var(--input) / <alpha-value>)',
				ring: 'oklch(var(--ring) / <alpha-value>)',
				background: 'oklch(var(--background) / <alpha-value>)',
				foreground: 'oklch(var(--foreground) / <alpha-value>)',
				// Surface colors
				surface: {
					DEFAULT: 'oklch(var(--surface) / <alpha-value>)',
					elevated: 'oklch(var(--surface-elevated) / <alpha-value>)',
					muted: 'oklch(var(--surface-muted) / <alpha-value>)'
				},
				// Glass colors
				glass: {
					bg: 'oklch(var(--glass-bg) / <alpha-value>)',
					border: 'oklch(var(--glass-border) / <alpha-value>)',
					shadow: 'oklch(var(--glass-shadow) / <alpha-value>)'
				},
				// Semantic colors
				primary: {
					DEFAULT: 'oklch(var(--primary) / <alpha-value>)',
					foreground: 'oklch(var(--primary-foreground) / <alpha-value>)'
				},
				secondary: {
					DEFAULT: 'oklch(var(--secondary) / <alpha-value>)',
					foreground: 'oklch(var(--secondary-foreground) / <alpha-value>)'
				},
				destructive: {
					DEFAULT: 'oklch(var(--destructive) / <alpha-value>)',
					foreground: 'oklch(var(--destructive-foreground) / <alpha-value>)'
				},
				success: {
					DEFAULT: 'oklch(var(--success) / <alpha-value>)',
					foreground: 'oklch(var(--success-foreground) / <alpha-value>)'
				},
				warning: {
					DEFAULT: 'oklch(var(--warning) / <alpha-value>)',
					foreground: 'oklch(var(--warning-foreground) / <alpha-value>)'
				},
				info: {
					DEFAULT: 'oklch(var(--info) / <alpha-value>)',
					foreground: 'oklch(var(--info-foreground) / <alpha-value>)'
				},
				muted: {
					DEFAULT: 'oklch(var(--muted) / <alpha-value>)',
					foreground: 'oklch(var(--muted-foreground) / <alpha-value>)'
				},
				accent: {
					DEFAULT: 'oklch(var(--accent) / <alpha-value>)',
					foreground: 'oklch(var(--accent-foreground) / <alpha-value>)'
				},
				popover: {
					DEFAULT: 'oklch(var(--popover) / <alpha-value>)',
					foreground: 'oklch(var(--popover-foreground) / <alpha-value>)'
				},
				card: {
					DEFAULT: 'oklch(var(--card) / <alpha-value>)',
					foreground: 'oklch(var(--card-foreground) / <alpha-value>)'
				}
			},
			borderRadius: {
				'xs': 'var(--radius-xs)',
				'sm': 'var(--radius-sm)',
				DEFAULT: 'var(--radius)',
				'md': 'var(--radius-md)',
				'lg': 'var(--radius-lg)',
				'xl': 'var(--radius-xl)',
				'2xl': 'var(--radius-2xl)',
				'full': 'var(--radius-full)'
			},
			fontFamily: {
				sans: ['Inter', 'system-ui', ...defaultTheme.fontFamily.sans],
				mono: ['JetBrains Mono', 'Consolas', ...defaultTheme.fontFamily.mono]
			},
			fontSize: {
				'xs': ['0.75rem', { lineHeight: '1rem' }],
				'sm': ['0.875rem', { lineHeight: '1.25rem' }],
				'base': ['1rem', { lineHeight: '1.5rem' }],
				'lg': ['1.125rem', { lineHeight: '1.75rem' }],
				'xl': ['1.25rem', { lineHeight: '1.75rem' }],
				'2xl': ['1.5rem', { lineHeight: '2rem' }],
				'3xl': ['1.875rem', { lineHeight: '2.25rem' }],
				'4xl': ['2.25rem', { lineHeight: '2.5rem' }],
				'5xl': ['3rem', { lineHeight: '1' }],
				'6xl': ['3.75rem', { lineHeight: '1' }],
				'7xl': ['4.5rem', { lineHeight: '1' }],
				'8xl': ['6rem', { lineHeight: '1' }],
				'9xl': ['8rem', { lineHeight: '1' }]
			},
			spacing: {
				'18': '4.5rem',
				'88': '22rem',
				'128': '32rem'
			},
			animation: {
				'fade-in': 'fadeIn 0.5s ease-out',
				'slide-in': 'slideIn 0.3s ease-out',
				'scale-in': 'scaleIn 0.2s ease-out',
				'float': 'float 6s ease-in-out infinite',
				'pulse-slow': 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
				'shimmer': 'shimmer 2s linear infinite'
			},
			keyframes: {
				fadeIn: {
					'0%': { opacity: '0', transform: 'translateY(10px)' },
					'100%': { opacity: '1', transform: 'translateY(0)' }
				},
				slideIn: {
					'0%': { transform: 'translateX(-100%)' },
					'100%': { transform: 'translateX(0)' }
				},
				scaleIn: {
					'0%': { transform: 'scale(0.9)', opacity: '0' },
					'100%': { transform: 'scale(1)', opacity: '1' }
				},
				float: {
					'0%, 100%': { transform: 'translateY(0px)' },
					'50%': { transform: 'translateY(-10px)' }
				},
				shimmer: {
					'0%': { transform: 'translateX(-100%)' },
					'100%': { transform: 'translateX(100%)' }
				}
			},
			backdropBlur: {
				'xs': '2px',
				'glass': '20px',
				'strong': '40px'
			},
			boxShadow: {
				'glass': '0 4px 6px -1px oklch(0 0 0 / 0.1), 0 2px 4px -2px oklch(0 0 0 / 0.1), inset 0 1px 0 0 oklch(1 0 0 / 0.1)',
				'glass-dark': '0 4px 6px -1px oklch(0 0 0 / 0.3), 0 2px 4px -2px oklch(0 0 0 / 0.3), inset 0 1px 0 0 oklch(1 0 0 / 0.05)',
				'brand': '0 10px 15px -3px oklch(var(--brand-primary) / 0.2), 0 4px 6px -4px oklch(var(--brand-primary) / 0.1)',
				'brand-lg': '0 20px 25px -5px oklch(var(--brand-primary) / 0.3), 0 8px 10px -6px oklch(var(--brand-primary) / 0.2)',
				'glow': '0 0 20px oklch(var(--brand-primary) / 0.3)',
				'glow-lg': '0 0 30px oklch(var(--brand-primary) / 0.4), 0 0 60px oklch(var(--brand-primary) / 0.2)',
				'inner-light': 'inset 0 1px 0 0 oklch(1 0 0 / 0.1)',
				'charcoal': '0 10px 15px -3px oklch(var(--brand-secondary) / 0.3), 0 4px 6px -4px oklch(var(--brand-secondary) / 0.2)'
			}
		}
	}
};

export default config;
